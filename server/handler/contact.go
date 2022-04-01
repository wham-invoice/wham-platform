package handler

import (
	"github.com/wham-invoice/wham-platform/db"
	"github.com/wham-invoice/wham-platform/server/route"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

type NewContactRequest struct {
	FirstName         string `json:"first_name" binding:"required"`
	LastName          string `json:"last_name" binding:"required"`
	Phone             string `json:"phone" binding:"required"`
	Email             string `json:"email" binding:"required"`
	Company           string `json:"company"`
	AddressFirstLine  string `json:"address_first_line"`
	AddressSecondLine string `json:"address_second_line"`
	Suburb            string `json:"suburb"`
	Postcode          string `json:"postcode"`
	Country           string `json:"country"`
}

// Contact returns a contact by ID.
var Contact = route.Endpoint{
	Method:  "GET",
	Path:    "/contact/get/:contact_id",
	Prereqs: route.Prereqs(EnsureContact()),
	Do: func(c *gin.Context) (interface{}, error) {
		contact := MustContact(c)

		return &contact, nil
	},
}

// UserContacts returns all contacts for a user.
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

// DeleteContact deletes a contact by ID.
var DeleteContact = route.Endpoint{
	Method:  "DELETE",
	Path:    "/contact/delete/:contact_id",
	Prereqs: route.Prereqs(EnsureContact()),
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		contact := MustContact(c)

		if err := contact.Delete(ctx, app); err != nil {
			return nil, errors.Trace(err)
		}

		return nil, nil
	},
}

// NewContact creates a new contact for the user.
var NewContact = route.Endpoint{
	Method: "POST",
	Path:   "/contact/new",
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		var req NewContactRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errors.Annotate(err, "cannot bind request")
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

func contactFromRequest(
	req NewContactRequest,
	userID string,
) *db.Contact {
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
