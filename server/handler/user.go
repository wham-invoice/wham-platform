package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/server/route"
)

// UserSummary returns total invoice amount and paid amount for the user.
var UserSummary = route.Endpoint{
	Method: "GET",
	Path:   "/user/summary",
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		summary, err := user.Summary(ctx, app)
		if err != nil {
			return nil, errors.Trace(err)
		}

		// if no invoices found, return.
		if summary.InvoiceTotal == 0 {
			return nil, nil
		}

		return &summary, nil
	},
}
