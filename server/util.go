package server

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
)

func getDataBase(c *gin.Context) (*db.App, error) {
	dbApp, ok := c.MustGet(database_conn_gin_key).(*db.App)
	if !ok {
		return nil, errors.New("could not get database connection from context")
	}
	return dbApp, nil
}

func getUserID(c *gin.Context) (string, error) {
	userID, ok := c.MustGet(user_id_gin_key).(string)
	if !ok {
		return "", errors.New("could not get user id from context")
	}
	return userID, nil
}
