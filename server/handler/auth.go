package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"
	"github.com/rstorr/wham-platform/util"

	"github.com/juju/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

type AuthRequest struct {
	UID     string `json:"uid" binding:"required"`      // UID from firebase authentication
	Code    string `json:"code" binding:"required"`     // server code from google used to get access_tokens server side
	IdToken string `json:"id_token" binding:"required"` // ID_token used to get google user info
}

type GoogleToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
}

var Auth = route.Endpoint{
	Method: "POST",
	Path:   "/auth",
	Do: func(c *gin.Context) (interface{}, error) {
		app := MustApp(c)

		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return nil, errors.Annotate(err, "could not bind request")
		}

		ctx := context.Background()
		user, err := authenticate(ctx, app, req)
		if err != nil {
			return nil, errors.Annotate(err, "Error authenticating user")
		}

		if err = SetSession(c, user); err != nil {
			return nil, errors.Annotate(err, "Error setting session")
		}

		return user, nil
	},
}

// handleAuth retrieves a User from the request and starts a session.
// If no User exists we create one and add to DB.
func authenticate(ctx context.Context, app *db.App, req AuthRequest) (*db.User, error) {
	user := &db.User{}

	userInfo, err := unpackIdToken(ctx, req.IdToken)
	if err != nil {
		return user, errors.Trace(err)
	}

	// see if we already have a user
	user, err = app.GetUser(ctx, req.UID)
	if err != nil {
		return user, errors.Trace(err)
	}

	if user == nil {
		authToken, err := googleExchange(req.Code)
		if err != nil {
			return user, errors.Trace(err)
		}

		if err = app.CreateAndSaveUser(
			ctx,
			req.UID,
			userInfo,
			authToken,
		); err != nil {
			return user, errors.Trace(err)
		}

		user, err = app.GetUser(ctx, req.UID)
		if err != nil {
			return user, errors.Trace(err)
		}
	}
	// TODO don't return user.oauth in response
	return user, errors.Trace(err)
}

// googleExchange retrieves a google access_token,id_token and other info
// from the serverAuthCode.
func googleExchange(serverAuthCode string) (oauth2.Token, error) {
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

	return token, nil
}

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
