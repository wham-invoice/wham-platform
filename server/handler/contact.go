package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/server/route"
)

var Contact = route.Endpoint{
	Method:  "GET",
	Path:    "/contact/get/:contact_id",
	Prereqs: route.Prereqs(EnsureContact()),
	Do: func(c *gin.Context) (interface{}, error) {
		contact := MustContact(c)

		return &contact, nil
	},
}
