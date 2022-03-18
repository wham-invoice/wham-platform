package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/server/route"
)

var UserSummary = route.Endpoint{
	Method:  "GET",
	Path:    "/user/summary",
	Prereqs: route.Prereqs(EnsureContact()),
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		summary, err := user.Summary(ctx, app)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return &summary, nil
	},
}
