package handler

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"
	"github.com/rstorr/wham-platform/util"
)

const (
	appDBKey       = "server:app_db"
	dbInvoiceKey   = "server:invoice"
	dbContactKey   = "server:contact"
	userSessionKey = "session:user"
)

// SetAppDB returns middleware that stores the application database in the gin
// context.
func SetAppDB(appDB *db.App) gin.HandlerFunc {
	return func(c *gin.Context) { c.Set(appDBKey, appDB) }
}

// MustApp returns the application database or panics.
func MustApp(c *gin.Context) *db.App {
	return c.MustGet(appDBKey).(*db.App)
}

// MustApp returns the application database or panics.
func MustInvoice(c *gin.Context) *db.Invoice {
	return c.MustGet(dbInvoiceKey).(*db.Invoice)
}

func MustUser(c *gin.Context) *db.User {
	s := MustSession(c)
	user := s.Get(userSessionKey).(db.User)
	return &user
}

func MustContact(c *gin.Context) *db.Contact {
	return c.MustGet(dbContactKey).(*db.Contact)
}

func MustSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	util.Logger.Infof("session: %v", session.ID())
	return session
}

func SetSession(c *gin.Context, user *db.User) error {
	session := sessions.Default(c)
	session.Set(userSessionKey, user)

	return session.Save()
}

// repoAccess looks up a repository in the database.
func invoiceExists(c *gin.Context, invoiceID string) (*db.Invoice, error) {
	app := MustApp(c)

	invoice, err := app.Invoice(context.Background(), invoiceID)
	if err == db.InvoiceNotFound {
		return nil, route.NotFound
	}

	return invoice, nil
}

// InvoiceExists returns middleware that extracts the value of :invoice_id and sets it in
// the context.
func InvoiceExists() gin.HandlerFunc {

	// This is the real handler, but it's convenient to use real errors.
	getInvoice := func(c *gin.Context) (*db.Invoice, error) {
		var req struct {
			ID string `uri:"invoice_id" binding:"required"`
		}
		if c.ShouldBindUri(&req); req.ID == "" {
			return nil, route.NotFound
		}
		return invoiceExists(c, req.ID)
	}

	// Get the repo, if allowed, and update the context or abort.
	return func(c *gin.Context) {
		invoice, err := getInvoice(c)
		if err != nil {
			route.Abort(c, err)
		} else {
			c.Set(dbInvoiceKey, invoice)
		}
	}
}

// repoAccess looks up a repository in the database.
func contactExists(c *gin.Context, contactID string) (*db.Contact, error) {
	app := MustApp(c)

	contact, err := app.Contact(context.Background(), contactID)
	if err == db.InvoiceNotFound {
		return nil, route.NotFound
	}

	return contact, nil
}

// ContactExists returns middleware that extracts the value of :contact_id and sets it in
// the context.
// TODO this naming convention is a bit weird
func ContactExists() gin.HandlerFunc {

	// This is the real handler, but it's convenient to use real errors.
	getContact := func(c *gin.Context) (*db.Contact, error) {
		var req struct {
			ID string `uri:"contact_id" binding:"required"`
		}
		if c.ShouldBindUri(&req); req.ID == "" {
			return nil, route.NotFound
		}
		return contactExists(c, req.ID)
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
