package handler

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"
)

const (
	dbAppKey       = "server:app_db"
	dbInvoiceKey   = "server:invoice"
	dbContactKey   = "server:contact"
	dbUserKey      = "server:user"
	userSessionKey = "session:user"
)

func SessionSetUserID(c *gin.Context, id string) error {
	session := sessions.Default(c)
	session.Set(userSessionKey, id)

	return session.Save()
}

func SessionGetUserID(c *gin.Context) string {
	session := sessions.Default(c)
	id := session.Get(userSessionKey).(string)

	return id
}

// SetAppDB returns middleware that stores the application database in the gin
// context.
func SetAppDB(appDB *db.App) gin.HandlerFunc {
	return func(c *gin.Context) { c.Set(dbAppKey, appDB) }
}

func SetUser(c *gin.Context, user *db.User) gin.HandlerFunc {
	return func(c *gin.Context) { c.Set(dbUserKey, user) }
}

// MustApp returns the application database or panics.
func MustApp(c *gin.Context) *db.App {
	return c.MustGet(dbAppKey).(*db.App)
}

// MustApp returns the application database or panics.
func MustInvoice(c *gin.Context) *db.Invoice {
	return c.MustGet(dbInvoiceKey).(*db.Invoice)
}

func MustUser(c *gin.Context) *db.User {
	user := c.MustGet(dbUserKey).(db.User)
	return &user
}

func MustContact(c *gin.Context) *db.Contact {
	return c.MustGet(dbContactKey).(*db.Contact)
}

// EnsureUser returns middleware that extracts the user_id from the session and sets the corresponding user in the context.
func EnsureUser() gin.HandlerFunc {
	getUser := func(c *gin.Context) (*db.User, error) {
		userID := SessionGetUserID(c)
		if userID == "" {
			return nil, route.NotFound
		}

		user, err := MustApp(c).User(context.Background(), userID)
		if err == db.UserNotFound {
			return nil, route.NotFound
		}

		return user, nil
	}

	return func(c *gin.Context) {
		user, err := getUser(c)
		if err != nil {
			route.Abort(c, err)
		} else {
			SetUser(c, user)
		}
	}
}

// EnsureInvoice returns middleware that extracts the value of :invoice_id and sets it in
// the context.
func EnsureInvoice() gin.HandlerFunc {
	getInvoice := func(c *gin.Context) (*db.Invoice, error) {
		var req struct {
			ID string `uri:"invoice_id" binding:"required"`
		}
		if c.ShouldBindUri(&req); req.ID == "" {
			return nil, route.NotFound
		}
		app := MustApp(c)

		invoice, err := app.Invoice(context.Background(), req.ID)
		if err == db.InvoiceNotFound {
			return nil, route.NotFound
		}

		return invoice, nil
	}

	return func(c *gin.Context) {
		invoice, err := getInvoice(c)
		if err != nil {
			route.Abort(c, err)
		} else {
			c.Set(dbInvoiceKey, invoice)
		}
	}
}

// EnsureContact returns middleware that extracts the value of :contact_id and sets it in
// the context.
func EnsureContact() gin.HandlerFunc {
	getContact := func(c *gin.Context) (*db.Contact, error) {
		var req struct {
			ID string `uri:"contact_id" binding:"required"`
		}
		if c.ShouldBindUri(&req); req.ID == "" {
			return nil, route.NotFound
		}

		app := MustApp(c)

		contact, err := app.Contact(context.Background(), req.ID)
		if err == db.ContactNotFound {
			return nil, route.NotFound
		}

		return contact, nil
	}

	return func(c *gin.Context) {
		contact, err := getContact(c)
		if err != nil {
			route.Abort(c, err)
		} else {
			c.Set(dbContactKey, contact)
		}
	}
}
