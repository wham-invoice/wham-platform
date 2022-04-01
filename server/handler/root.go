// Package handler implements the external codelingo API.
//
// It also implements the internal admin API, which probably ought to be a
// separate service entirely, and not exposed to the outside world.
package handler

import (
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

type Session interface {
	GetUser(
		c *gin.Context,
		app *db.App,
	) (*db.User, error)
}

// Config configures an api server.
type Config struct {
	AllowOrigin string
	AppDB       *db.App
	RedisStore  *redis.Store
	Session     Session
}

const sessionName = "user_session"

// Validate returns an error if the Config is not sensible.
func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.AllowOrigin) == "" {
		return errors.New("blank AllowOrigin")
	}

	if cfg.AppDB == nil {
		return errors.New("missing AppDB")
	}

	if cfg.Session == nil {
		return errors.New("missing Session")
	}

	return nil
}

// Root returns the complete API.
func Root(cfg Config) (route.Installer, error) {

	auth, err := authorized(cfg)
	if err != nil {
		return nil, errors.Annotate(err, "cannot create authorized prereq")
	}

	root, err := rootPreReqs(cfg)
	if err != nil {
		return nil, errors.Annotate(err, "cannot create root prereq")
	}

	return route.Group{
		Path: "/",
		// Always handle panics, always log requests.
		Prereqs: root,
		Installers: route.Installers(
			// Trivial "is it running" health check.
			route.Endpoint{
				Method: "GET",
				Path:   "/",
				Do: func(c *gin.Context) (interface{}, error) {
					return nil, nil
				},
			},
			// paths that don't require auth.
			route.Group{
				Path: "/",
				Installers: route.Installers(
					Auth,
					ViewInvoice,
					PDF,
				),
			},
			// paths that require auth.
			route.Group{
				Path:    "/",
				Prereqs: auth,
				Installers: route.Installers(
					Invoice,
					EmailInvoice,
					NewInvoice,
					DeleteInvoice,
					UserInvoices,
					Contact,
					UserContacts,
					NewContact,
					DeleteContact,
					UserSummary,
				),
			},
		),
	}, nil
}

// authorized returns a prereq that checks for a valid user.
func authorized(cfg Config) ([]gin.HandlerFunc, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Annotate(err, "bad config")
	}

	return route.Prereqs(
		SetSession(cfg.Session),
		EnsureUser(),
	), nil
}

// rootPreReqs returns the root prereqs.
func rootPreReqs(cfg Config) ([]gin.HandlerFunc, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Annotate(err, "bad config")
	}

	return route.Prereqs(
		gin.Recovery(),
		gin.Logger(),
		sessions.Sessions(sessionName, *cfg.RedisStore),
		setUpCors(cfg),
		SetAppDB(cfg.AppDB),
	), nil
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
