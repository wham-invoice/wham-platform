package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/wham-invoice/wham-platform/db"
	"github.com/wham-invoice/wham-platform/server/route"
)

type RealSession struct{}

// GetUser returns the user from the session.
func (RealSession) GetUser(
	c *gin.Context,
	app *db.App,
) (*db.User, error) {

	userID := SessionGetUserID(c)
	if userID == "" {
		return nil, route.NotFound
	}

	user, err := MustApp(c).User(context.Background(), userID)
	if err == db.UserNotFound {
		return nil, route.NotFound
	}

	return user, nil

}
