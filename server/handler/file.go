package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/server/route"
)

// PDF returns the invoice pdf file by id
var PDF = route.Endpoint{
	Method: "GET",
	Path:   "/pdf/:pdf_id",
	Do: func(c *gin.Context) (interface{}, error) {

		app := MustApp(c)

		var req struct {
			ID string `uri:"pdf_id" binding:"required"`
		}
		if c.ShouldBindUri(&req); req.ID == "" {
			return nil, route.NotFound
		}

		fileName := fmt.Sprintf("%s", req.ID)
		body, err := app.PDF(c.Request.Context(), fileName)
		if err != nil {
			return nil, errors.Trace(err)
		}

		// ISTM we're pulling the pdf body into memory then sending to the client.
		// TODO what happens to PDF after we've sent it?
		cType := http.DetectContentType(body)
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileName))
		c.Data(http.StatusOK, cType, body)

		return nil, nil
	},
}
