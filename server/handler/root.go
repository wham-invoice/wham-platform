// Package handler implements the external codelingo API.
//
// It also implements the internal admin API, which probably ought to be a
// separate service entirely, and not exposed to the outside world.
package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"
)

// Config configures an api server.
type Config struct {
	AllowOrigin string
	AppDB       *db.App
	RedisStore  *redis.Store
}

// Validate returns an error if the Config is not sensible.
func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.AllowOrigin) == "" {
		return errors.New("blank AllowOrigin")
	}

	if cfg.AppDB == nil {
		return errors.New("missing AppDB")
	}

	return nil
}

// Root returns the complete API.
func Root(cfg Config) (route.Installer, error) {

	auth, err := authorized(cfg)
	if err != nil {
		return nil, errors.Annotate(err, "cannot create authorized prereq")
	}

	unauth, err := unauthorized(cfg)
	if err != nil {
		return nil, errors.Annotate(err, "cannot create unauthorized prereq")
	}

	return route.Group{
		Path: "/",
		// Always handle panics, always log requests.
		Prereqs: route.Prereqs(gin.Recovery(), gin.Logger()),
		Installers: route.Installers(
			// Trivial "is it running" health check.
			route.Endpoint{
				Method: "GET",
				Path:   "/",
				Do: func(c *gin.Context) (interface{}, error) {
					return nil, nil
				},
			},
			route.Group{
				Path:    "/",
				Prereqs: unauth,
				Installers: route.Installers(
					Auth,
					Invoice,
				//	PDF,
				),
			},
			route.Group{
				Path:    "/",
				Prereqs: auth,
				Installers: route.Installers(
					EmailInvoice,
					NewInvoice,
					AllInvoices,
					Contact,
				),
			},
		),
	}, nil
}

func authorized(cfg Config) ([]gin.HandlerFunc, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Annotate(err, "bad config")
	}

	return route.Prereqs(setUpCors(cfg), SetAppDB(cfg.AppDB), user()), nil
}

func unauthorized(cfg Config) ([]gin.HandlerFunc, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Annotate(err, "bad config")
	}

	return route.Prereqs(setUpCors(cfg), SetAppDB(cfg.AppDB)), nil
}

func user() gin.HandlerFunc {
	return func(c *gin.Context) {
		app := MustApp(c)
		session := sessions.Default(c)
		userID := session.Get(userSessionKey)
		if userID == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "unauthorized",
			})
			route.Abort(c, errors.New("cannot get user ID from session"))
		}

		user, err := app.GetUser(context.Background(), userID.(string))
		if err != nil {
			route.Abort(c, errors.Annotate(err, "cannot get user"))
		}

		c.Set(dbUserKey, user)
		c.Next()
	}
}

func setUpCors(cfg Config) gin.HandlerFunc {
	// CORS for everything.
	return cors.New(cors.Config{
		AllowOrigins: []string{
			cfg.AllowOrigin,
		},
		AllowMethods: []string{
			http.MethodDelete,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"Content-Length",
			"Origin",
		},
	})
}
