package server

import (
	"encoding/gob"
	"net/http"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/handler"
	"golang.org/x/oauth2"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

func init() {
	gob.Register(oauth2.Token{})
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(dbApp *db.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", dbApp)
		c.Next()
	}
}

//AuthMiddleware ensures that a current session exists.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
		}
		c.Set("userID", userID)
		c.Next()
	}
}

func Run(dbApp *db.App) error {
	// TODO i think 'secret' needs to be an actual secret...
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		return errors.Trace(err)
	}

	router := gin.Default()
	router.Use(sessions.Sessions("user_session", store))
	router.Use(ApiMiddleware(dbApp))
	router.POST("/auth", handler.Authenticate)

	auth := router.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/invoice/:id/pdf", handler.GenerateInvoicePDFHandler)
		auth.POST("/invoice/email", handler.EmailInvoice)
		auth.POST("/invoice/new", handler.EmailInvoice)
	}

	return router.Run(":8080")
}
