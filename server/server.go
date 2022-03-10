package server

import (
	"context"
	"encoding/gob"
	"fmt"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/handler"
	"golang.org/x/oauth2"

	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

const user_id_gin_key = "user_id_key"
const user_id_session_key = "user_id_session_key"
const database_conn_gin_key = "database_connection_key"

func init() {
	gob.Register(oauth2.Token{})
	gob.Register(db.User{})
}

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(dbApp *db.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(database_conn_gin_key, dbApp)
		c.Next()
	}
}

// TODO requests need to run async.
// TODO look at platform/service/api/main.go for reference.
func Run(ctx context.Context) error {
	var cfg handler.Config
	addr, err := configure(ctx, &cfg)
	if err != nil {
		return errors.Annotate(err, "cannot configure")
	}

	root, err := handler.Root(cfg)
	if err != nil {
		return errors.Annotate(err, "cannot create handler")
	}

	ngin := gin.New()
	root.Install(&ngin.RouterGroup)
	return ngin.Run(addr)
}

func configure(ctx context.Context, cfg *handler.Config) (string, error) {

	addr := "0.0.0.0:8080"

	cfg.AllowOrigin = "http://localhost:3000"

	// TODO i think 'secret' needs to be an actual secret...
	store, err := redis.NewStore(10, "tcp", fmt.Sprintf("%s:%d", "localhost", 6379), "", []byte("secret"))
	if err != nil {
		return "", errors.Annotate(err, "cannot set up redis store")
	}
	cfg.RedisStore = &store

	// Set this up last, once everything else looks like it worked.
	// Don't bother to close, it should live as long as the process anyway.
	cfg.AppDB, err = db.Init(ctx)
	if err != nil {
		return "", errors.Annotate(err, "cannot set up application DB")
	}

	return addr, nil
}
