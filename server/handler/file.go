package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/server/route"
)

var PDF = route.Endpoint{
	Method: "GET",
	Path:   "/pdf/:pdf_id",
	Do: func(c *gin.Context) (interface{}, error) {
		// TODO grab file from db. Store in memory/disk? on server. then return it?

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

		cType := http.DetectContentType(body)
		c.Header("Access-Control-Expose-Headers", "Content-Disposition")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.pdf", fileName))
		c.Data(http.StatusOK, cType, body)

		return nil, nil
	},
}
