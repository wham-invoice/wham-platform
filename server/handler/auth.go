package handler

import (
	"context"
	"net/http"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
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

func Authenticate(c *gin.Context) {
	//TODO if request does not have 'x-requested-with' header this could be a CSRF
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	dbApp, ok := c.MustGet("databaseConn").(*db.App)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("could not get database connection"))
	}

	user, err := handleAuth(c, dbApp, req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.IndentedJSON(http.StatusOK, user)
}

// handleAuth retrieves a User from the request and starts a session.
// If no User exists we create one and add to DB.
func handleAuth(c *gin.Context, app *db.App, req AuthRequest) (*db.User, error) {
	ctx := context.Background()
	user := &db.User{}

	userInfo, err := auth.UnpackIdToken(ctx, req.IdToken)
	if err != nil {
		return user, errors.Trace(err)
	}

	// see if we already have a user
	user, err = app.GetUser(ctx, &db.UserRequest{ID: req.UID})
	if err != nil {
		return user, errors.Trace(err)
	}

	if user == nil {
		authToken, err := auth.GoogleExchange(req.Code)
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

		user, err = app.GetUser(ctx, &db.UserRequest{ID: req.UID})
		if err != nil {
			return user, errors.Trace(err)
		}

		s := sessions.Default(c)
		if err = setSession(s, user.ID); err != nil {
			return user, errors.Trace(err)
		}
	}

	return user, nil
}

func setSession(session sessions.Session, id string) error {
	session.Set("user_id", id)

	return session.Save()
}
