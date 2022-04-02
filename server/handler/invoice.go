package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wham-invoice/wham-platform/db"
	"github.com/wham-invoice/wham-platform/email"
	"github.com/wham-invoice/wham-platform/pdf"
	"github.com/wham-invoice/wham-platform/server/route"

	"github.com/juju/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type EmailInvoiceRequest struct {
	ID string `json:"invoice_id" binding:"required"`
}

type NewInvoiceRequest struct {
	ContactID   string  `json:"contact_id" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Hours       float32 `json:"hours" binding:"required"`
	Rate        float32 `json:"rate" binding:"required"`
	DueDate     string  `json:"due_date" binding:"required"`
}

// Invoice returns the invoice by id
var Invoice = route.Endpoint{
	Method:  "GET",
	Path:    "/invoice/get/:invoice_id",
	Prereqs: route.Prereqs(EnsureInvoice()),
	Do: func(c *gin.Context) (interface{}, error) {
		invoice := MustInvoice(c)

		return &invoice, nil
	},
}

// ViewInvoice is a handler for viewing an invoice on wham-web. We return the associated invoice,
// user and contact info.
var ViewInvoice = route.Endpoint{
	Method:  "GET",
	Path:    "/invoice/view/:invoice_id",
	Prereqs: route.Prereqs(EnsureInvoice()),
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		invoice := MustInvoice(c)

		// ISTM invoice detail just munges invoice and its contact into one response.
		detail, err := invoice.Detail(ctx, app)
		return &detail, err
	},
}

// DeleteInvoice is a handler for deleting an invoice by its ID.
var DeleteInvoice = route.Endpoint{
	Method:  "DELETE",
	Path:    "/invoice/delete/:invoice_id",
	Prereqs: route.Prereqs(EnsureInvoice()),
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		invoice := MustInvoice(c)

		if err := invoice.Delete(ctx, app); err != nil {
			return nil, errors.Trace(err)
		}

		return nil, nil
	},
}

// TODO pagination
// UserInvoices returns all invoices for a user.
var UserInvoices = route.Endpoint{
	Method: "GET",
	Path:   "/user/invoices",
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		invoices, err := user.Invoices(ctx, app)
		if err != nil {
			return nil, errors.Annotate(err, "cannot get invoices")
		}

		return invoices, nil
	},
}

// NewInvoice is a handler for creating a new invoice.
var NewInvoice = route.Endpoint{
	Method: "POST",
	Path:   "/invoice/new",
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		var req NewInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			return nil, errors.Annotate(err, "cannot bind request")
		}

		newInvoice, err := invoiceFromRequest(req, user.ID)
		if err != nil {
			return nil, errors.Annotate(err, "cannot create invoice from request")
		}

		contact, err := app.Contact(ctx, newInvoice.ContactID)
		if err != nil {
			return nil, errors.Annotate(err, "cannot get contact ")
		}

		pdfBuilder := &pdf.Builder{
			App:     app,
			Invoice: newInvoice,
			User:    user,
			Contact: contact}
		pdfID, err := pdf.CreatePDF(ctx, *pdfBuilder)
		if err != nil {
			return nil, errors.Annotate(err, "cannot create PDF from invoice")
		}

		newInvoice.PDFID = pdfID

		id, err := app.AddInvoice(ctx, newInvoice)
		if err != nil {
			return nil, errors.Annotate(err, "cannot add new invoice")
		}

		invoice, err := app.Invoice(ctx, id)
		if err != nil {
			return nil, errors.Annotate(err, "cannot get new invoice")
		}

		return invoice, nil
	},
}

// TODO invoice_id should be in path then use MustInvoice.
var EmailInvoice = route.Endpoint{
	Method: "POST",
	Path:   "/invoice/email",
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		var req EmailInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Annotate(err, "cannot bind request"))
			return nil, nil
		}

		invoice, err := app.Invoice(ctx, req.ID)
		if err != nil {
			return nil, errors.Trace(err)
		}

		contact, err := app.Contact(ctx, invoice.ContactID)
		if err != nil {
			return nil, errors.Trace(err)
		}

		if err := emailInvoice(ctx, invoice, user, contact); err != nil {
			return nil, errors.Trace(err)
		}

		return nil, nil
	},
}

// TODO config should be stored in config file.
func emailInvoice(
	ctx context.Context,
	invoice *db.Invoice,
	user *db.User,
	contact *db.Contact,
) error {
	b, err := ioutil.ReadFile("./secrets/google_web_client_credentials.json")
	if err != nil {
		return errors.Trace(err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope, gmail.GmailSendScope)
	if err != nil {
		return errors.Trace(err)
	}

	httpClient := config.Client(context.Background(), &user.OAuth)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return errors.Trace(err)
	}

	invoiceURL := fmt.Sprintf("http://localhost:3000/invoice/%s", invoice.ID)
	body := fmt.Sprintf("Hi %s,\n\n"+
		"Your invoice is ready.\n\n"+
		"To view and download it please visit: %s "+
		"Thanks.\n"+
		"%s", contact.FirstName, invoiceURL, user.FirstName)

	return errors.Trace(
		email.GmailSend(service, "me", contact.Email, "Invoice", body),
	)
}

func invoiceFromRequest(req NewInvoiceRequest, userID string) (*db.Invoice, error) {
	dueDate, err := time.Parse("2006-01-02T00:00:00.000", req.DueDate)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &db.Invoice{
		UserID:      userID,
		ContactID:   req.ContactID,
		Description: req.Description,
		Rate:        req.Rate,
		Hours:       req.Hours,
		IssueDate:   time.Now(),
		DueDate:     dueDate,
	}, nil
}
