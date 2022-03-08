package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/email"
	"github.com/rstorr/wham-platform/pdf"
	"github.com/rstorr/wham-platform/server/route"
	"github.com/rstorr/wham-platform/util"

	"github.com/juju/errors"
	uuid "github.com/satori/go.uuid"
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
	Prereqs: route.Prereqs(InvoiceAccess()),
	Do: func(c *gin.Context) (interface{}, error) {
		invoice := MustInvoice(c)

		return &invoice, nil
	},
}

// TODO pagination
var AllInvoices = route.Endpoint{
	Method:  "GET",
	Path:    "/invoice/getAll",
	Prereqs: route.Prereqs(InvoiceAccess()),
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := context.Background()
		app := MustApp(c)
		user := MustUser(c)

		invoices, err := user.Invoices(ctx, app)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Annotate(err, "cannot get invoices"))
		}

		return invoices, nil
	},
}

// NOTE: shouldn't this return the created invoice?
var NewInvoice = route.Endpoint{
	Method:  "POST",
	Path:    "/invoice/new",
	Prereqs: route.Prereqs(InvoiceAccess()),
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := context.Background()
		app := MustApp(c)
		user := MustUser(c)

		var req NewInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		i := invoiceFromRequest(req, user.ID)
		if err := newInvoice(ctx, app, i); err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Annotate(err, "cannot add new invoice"))
		}

		return nil, nil
	},
}

var EmailInvoice = route.Endpoint{
	Method:  "GET",
	Path:    "/invoice/get/:invoice_id",
	Prereqs: route.Prereqs(InvoiceAccess()),
	Do: func(c *gin.Context) (interface{}, error) {
		ctx := context.Background()
		invoice := MustInvoice(c)
		app := MustApp(c)
		user := MustUser(c)

		var req EmailInvoiceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		contact, err := app.GetContact(ctx, invoice.ContactID)
		if err != nil {
			util.Logger.Error(errors.ErrorStack(err))
			c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
		}

		if err := emailInvoice(ctx, invoice, user, contact); err != nil {
			c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
		}

		return nil, nil
	},
}

func newInvoice(ctx context.Context, dbApp *db.App, invoice *db.Invoice) error {

	pdfID, err := createPDF(ctx, dbApp, invoice)
	if err != nil {
		return errors.Trace(err)
	}
	invoice.PDFID = pdfID

	if err := dbApp.AddInvoice(ctx, invoice); err != nil {
		return errors.Trace(err)
	}

	return nil
}

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

	return errors.Trace(email.GmailSend(service, "me", contact.Email, "Invoice", body))
}

// createPDF creates a PDF from an invoice ID and stores the file in firebase.
// We delete the file from local disk. Finally we return the ID of the file in firebase.
func createPDF(ctx context.Context, dbApp *db.App, invoice *db.Invoice) (string, error) {

	pdfID := uuid.NewV4().String()

	filePath := fmt.Sprintf("invoices/%s.pdf", pdfID)

	if err := pdf.Construct(pdf.PDFConstructor{
		Invoice:    invoice,
		OutputPath: filePath,
	}); err != nil {
		return "", errors.Trace(err)
	}
	if err := dbApp.UploadPDFDeleteLocal(ctx, pdfID, filePath); err != nil {
		return "", errors.Trace(err)
	}
	if err := os.Remove(filePath); err != nil {
		return "", errors.Trace(err)
	}

	return pdfID, nil
}

func invoiceFromRequest(req NewInvoiceRequest, userID string) *db.Invoice {
	return &db.Invoice{
		UserID:      userID,
		Description: req.Description,
		Rate:        req.Rate,
		Hours:       req.Hours,
	}
}
