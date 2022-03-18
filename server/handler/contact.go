package handler

import (
	"github.com/rstorr/wham-platform/server/route"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
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

var UserContacts = route.Endpoint{
	Method:  "GET",
	Path:    "/user/contacts",
	Prereqs: route.Prereqs(EnsureContact()),
	Do: func(c *gin.Context) (interface{}, error) {
		user := MustUser(c)
		contacts, err := user.Contacts(c.Request.Context(), MustApp(c))
		if err != nil {
			return nil, errors.Trace(err)
		}

		return contacts, nil
	},
}
