package server

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/util"
	"golang.org/x/oauth2"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

const user_id_gin_key = "user_id_key"
const user_id_session_key = "user_id_session_key"
const database_conn_gin_key = "database_connection_key"

func init() {
	gob.Register(oauth2.Token{})
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(dbApp *db.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(database_conn_gin_key, dbApp)
		c.Next()
	}
}

// AuthMiddleware ensures that a current session exists.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(user_id_session_key)
		if userID == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
		}
		c.Set(user_id_gin_key, userID)
		c.Next()
	}
}

// TODO requests need to run async.
func Run(dbApp *db.App) error {
	// TODO i think 'secret' needs to be an actual secret...
	storeAddr := fmt.Sprintf("%s:%d", "localhost", 6379)
	store, err := redis.NewStore(10, "tcp", storeAddr, "", []byte("secret"))
	if err != nil {
		return errors.Trace(err)
	}
	util.Logger.Infof("connected to redis on %s", storeAddr)

	router := gin.Default()
	// NOTE does multiple sessions work?
	router.Use(sessions.Sessions("user_session", store))
	router.Use(ApiMiddleware(dbApp))
	router.POST("/auth", authenticateHandler)

	auth := router.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.POST("/invoice/email", emailInvoiceHandler)
		auth.POST("/invoice/new", newInvoiceHandler)
		auth.POST("/invoice/get", getInvoiceHandler)
	}

	return router.Run(":8080")
}
