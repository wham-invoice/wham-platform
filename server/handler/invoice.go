package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/email"
	"github.com/rstorr/wham-platform/pdf"
	"github.com/rstorr/wham-platform/server/route"

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
}

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

		return invoice.Detail(ctx, app)
	},
}

// TODO pagination
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

var NewInvoice = route.Endpoint{
	Method: "POST",
	Path:   "/invoice/new",
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := c.Request.Context()
		app := MustApp(c)
		user := MustUser(c)

		var req NewInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.Annotate(err, "cannot bind request"))
			return nil, nil
		}

		newInvoice := invoiceFromRequest(req, user.ID)

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
	b, err := ioutil.ReadFile("/Users/work/go/src/github.com/rstorr/wham-platform/secrets/google_web_client_credentials.json")
	if err != nil {
		return errors.Trace(err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope)
	if err != nil {
		return errors.Trace(err)
	}

	httpClient := config.Client(context.Background(), &user.OAuth)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return errors.Trace(err)
	}

	invoiceURL := fmt.Sprintf("https://whaminvoice.co.nz/invoice/%s", invoice.PDFID)
	body := fmt.Sprintf("Hi %s,\n\n"+
		"Your invoice is ready.\n\n"+
		"To view and download it please visit: %s "+
		"Thanks.\n"+
		"%s", contact.FirstName, invoiceURL, user.FirstName)

	return errors.Trace(
		email.GmailSend(service, "me", contact.Email, "Invoice", body),
	)
}

func invoiceFromRequest(req NewInvoiceRequest, userID string) *db.Invoice {
	return &db.Invoice{
		UserID:      userID,
		ContactID:   req.ContactID,
		Description: req.Description,
		Rate:        req.Rate,
		Hours:       req.Hours,
	}
}
