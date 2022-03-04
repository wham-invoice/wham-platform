package server

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

type GetInvoiceRequest struct {
	ID string `json:"invoice_id" binding:"required"`
}

func newInvoiceHandler(c *gin.Context) {
	ctx := context.Background()

	var req NewInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userID, err := getUserID(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	dbApp, err := getDataBase(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	i := invoiceFromRequest(req, userID)
	if err := NewInvoice(ctx, dbApp, i); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	c.Status(http.StatusOK)
}

func getInvoiceHandler(c *gin.Context) {
	ctx := context.Background()

	var req GetInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	dbApp, err := getDataBase(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	invoice, err := dbApp.GetInvoice(ctx, &db.InvoiceRequest{ID: req.ID})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

func emailInvoiceHandler(c *gin.Context) {
	ctx := context.Background()

	var req EmailInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.Logger.Error(errors.ErrorStack(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	dbApp, err := getDataBase(c)
	if err != nil {
		util.Logger.Error(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	userID, err := getUserID(c)
	if err != nil {
		util.Logger.Error(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	user, err := dbApp.GetUser(ctx, &db.UserRequest{ID: userID})
	if err != nil {
		util.Logger.Error(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	i, err := dbApp.GetInvoice(ctx, &db.InvoiceRequest{ID: req.ID})
	if err != nil {
		util.Logger.Error(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	contact, err := dbApp.GetContact(ctx, &db.ContactRequest{ID: i.ContactID})
	if err != nil {
		util.Logger.Error(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	if err := emailInvoice(ctx, i, user, contact); err != nil {
		util.Logger.Error(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}
}

func NewInvoice(ctx context.Context, dbApp *db.App, invoice *db.Invoice) error {

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

	util.Logger.Infof("%v", &user.OAuth)

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
