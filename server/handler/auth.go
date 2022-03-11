package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"
	"github.com/rstorr/wham-platform/util"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type AuthRequest struct {
	UID     string `json:"uid" binding:"required"`      // UID from firebase authentication
	IdToken string `json:"id_token" binding:"required"` // ID_token used to get google user info
	Code    string `json:"code" binding:"required"`     // server code from google used to get access_tokens server side
}

type GoogleToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
}

// TODO we need to keep an eye on oauth2 expiries and refresh tokens when necessary.
var Auth = route.Endpoint{
	Method: "POST",
	Path:   "/auth",
	Do: func(c *gin.Context) (interface{}, error) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errors.Annotate(err, "could not bind request")
		}

		app := MustApp(c)

		user, err := authenticate(c.Request.Context(), app, req)
		if err != nil {
			return nil, errors.Annotate(err, "Error authenticating user")
		}

		if err = SessionSetUserID(c, user.ID); err != nil {
			return nil, errors.Annotate(err, "Error setting session")
		}

		return user.Sanitize(), nil
	},
}

// authenticate retrieves a User from the request. If no User exists we create one and add to DB.
func authenticate(ctx context.Context, app *db.App, req AuthRequest) (*db.User, error) {

	// Get info on what user is logging in
	userInfo, err := unpackIdToken(ctx, req.IdToken)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// see if we already have a user
	user, err := app.User(ctx, req.UID)
	if err != nil {
		return user, errors.Trace(err)
	}

	// if we have a user return it
	if user != nil {
		return user, nil
	}

	// if we don't have a user create one.

	// get an oAuth token from google
	authToken, err := tokenFromGoogle(req.Code)
	if err != nil {
		return user, errors.Trace(err)
	}

	// create a new user and add to DB
	if err = app.NewUser(
		ctx,
		req.UID,
		userInfo,
		authToken,
	); err != nil {
		return user, errors.Trace(err)
	}

	// get new user
	user, err = app.User(ctx, req.UID)
	if err != nil {
		return user, errors.Trace(err)
	}

	return user, nil
}

// tokenFromGoogle sends the serverAuthCode to google to get an oauth2 token.
// We manually set the expiry time to be 55 mins from now. This is because the token returned from
// google does not have the expiry set correctly.
// TODO: this config should be moved to config file
func tokenFromGoogle(serverAuthCode string) (oauth2.Token, error) {
	var token oauth2.Token

	v := url.Values{
		"Content-Type":  {"application/x-www-form-urlencoded; charset=utf-8"},
		"code":          {serverAuthCode},
		"client_id":     {util.GetConfigValue("oauth2.client_id")},
		"client_secret": {util.GetConfigValue("oauth2.client_secret")},
		"redirect_uri":  {"https://wham-ad61b.firebaseapp.com/__/auth/handler"},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", v)
	if err != nil {
		return token, errors.Trace(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return token, errors.Errorf("bad request error hitting google auth.")
	}

	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return token, errors.Trace(err)
	}

	now := time.Now()
	expiry := now.Add(time.Duration(55) * time.Minute)
	token.Expiry = expiry

	return token, nil
}

// unpackIdToken takes an id_token unmarshals it into a UserInfo.
func unpackIdToken(ctx context.Context, token string) (db.UserInfo, error) {
	var info db.UserInfo

	payload, err := idtoken.Validate(ctx, token,
		util.GetConfigValue("oauth2.client_id"))
	if err != nil {
		return info, errors.Trace(err)
	}

	// Convert map to json string
	jsonStr, err := json.Marshal(payload.Claims)
	if err != nil {
		return info, errors.Trace(err)
	}

	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, &info); err != nil {
		return info, errors.Trace(err)
	}

	return info, nil
}
