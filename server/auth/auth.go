package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rstorr/wham-platform/util"

	"github.com/juju/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

//googleExchange retrieves a google access_token,id_token and other info
// from the serverAuthCode.
func GoogleExchange(serverAuthCode string) (oauth2.Token, error) {
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

	fmt.Println(resp.Body)

	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return token, errors.Trace(err)
	}

	return token, nil
}

type UserInfo struct {
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Name       string `json:"name"`
	Sub        string `json: "sub"`
	// An identifier for the user,
	// unique among all Google accounts and never reused.
	// A Google account can have multiple email addresses at different points in time, but the sub value is never changed.
	// Use sub within your application as the unique-identifier key for the user
}

func UnpackIdToken(ctx context.Context, token string) (UserInfo, error) {
	var info UserInfo

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
