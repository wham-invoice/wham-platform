package handler

import (
	"context"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/route"
)

const (
	appDBKey       = "server:app_db"
	dbInvoiceKey   = "server:invoice"
	dbUserKey      = "server:user"
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
	return c.MustGet(dbUserKey).(*db.User)
}

func MustContact(c *gin.Context) *db.Contact {
	return c.MustGet(dbContactKey).(*db.Contact)
}

func SetSession(session sessions.Session, user *db.User) error {
	session.Set(userSessionKey, user)

	return session.Save()
}

// repoAccess looks up a repository in the database.
func invoiceAccess(c *gin.Context, invoiceID string) (*db.Invoice, error) {
	app := MustApp(c)

	invoice, err := app.GetInvoice(context.Background(), &db.InvoiceRequest{ID: invoiceID})
	if err == db.InvoiceNotFound {
		return nil, route.NotFound
	}

	return invoice, nil
}

// InvoiceAccess returns middleware that extracts the value of :invoice_id and sets it in
// the context.
func InvoiceAccess() gin.HandlerFunc {

	// This is the real handler, but it's convenient to use real errors.
	getInvoice := func(c *gin.Context) (*db.Invoice, error) {
		var i struct {
			ID string `uri:"invoice_id" binding:"required"`
		}
		if c.ShouldBindUri(&i); i.ID == "" {
			return nil, route.NotFound
		}
		return invoiceAccess(c, i.ID)
	}

	// Get the repo, if allowed, and update the context or abort.
	return func(c *gin.Context) {
		repo, err := getInvoice(c)
		if err != nil {
			route.Abort(c, err)
		} else {
			c.Set(dbInvoiceKey, repo)
		}
	}
}

// repoAccess looks up a repository in the database.
func contactAccess(c *gin.Context, invoiceID string) (*db.Invoice, error) {
	app := MustApp(c)

	invoice, err := app.GetInvoice(context.Background(), &db.InvoiceRequest{ID: invoiceID})
	if err == db.InvoiceNotFound {
		return nil, route.NotFound
	}

	return invoice, nil
}

// InvoiceAccess returns middleware that extracts the value of :invoice_id and sets it in
// the context.
func ContactAccess() gin.HandlerFunc {

	// This is the real handler, but it's convenient to use real errors.
	getContact := func(c *gin.Context) (*db.Invoice, error) {
		var contact struct {
			ID string `uri:"contact_id" binding:"required"`
		}
		if c.ShouldBindUri(&contact); contact.ID == "" {
			return nil, route.NotFound
		}
		return contactAccess(c, contact.ID)
	}

	return func(c *gin.Context) {
		repo, err := getContact(c)
		if err != nil {
			route.Abort(c, err)
		} else {
			c.Set(dbContactKey, repo)
		}
	}
}
