package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"
)

type RealSession struct{}

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
