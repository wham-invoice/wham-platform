package handler

import (
	"net/http"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

type NewContactRequest struct {
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	Phone             string `json:"phone" binding:"required"`
	Email             string `json:"email" binding:"required"`
	Company           string `json:"company" binding:"required"`
	AddressFirstLine  string `json:"address_first_line" binding:"required"`
	AddressSecondLine string `json:"address_second_line" binding:"required"`
	Suburb            string `json:"suburb" binding:"required"`
	Postcode          string `json:"postcode" binding:"required"`
	Country           string `json:"country" binding:"required"`
}

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
	Method: "GET",
	Path:   "/user/contacts",
	Do: func(c *gin.Context) (interface{}, error) {
		user := MustUser(c)
		contacts, err := user.Contacts(c.Request.Context(), MustApp(c))
		if err != nil {
			return nil, errors.Trace(err)
		}

		return contacts, nil
	},
}

var NewContact = route.Endpoint{
	Method: "POST",
	Path:   "/contact/new",
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		var req NewContactRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Annotate(err, "cannot bind request"))
			return nil, nil
		}

		newContact := contactFromRequest(req, user.ID)

		id, err := app.AddContact(ctx, newContact)
		if err != nil {
			return nil, errors.Annotate(err, "cannot add new contact")
		}

		contact, err := app.Contact(ctx, id)
		if err != nil {
			return nil, errors.Annotate(err, "cannot get new contact")
		}

		return contact, nil
	},
}

func contactFromRequest(req NewContactRequest, userID string) *db.Contact {

	return &db.Contact{
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,
		Company:   req.Company,
		Address: &db.Address{
			FirstLine:  req.AddressFirstLine,
			SecondLine: req.AddressSecondLine,
			Suburb:     req.Suburb,
			Postcode:   req.Postcode,
			Country:    req.Country,
		},
	}
}
