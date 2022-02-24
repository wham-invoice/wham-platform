package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/pdf"
	"github.com/rstorr/wham-platform/server/email"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

// TODO - this should also create the invoice.
// GenerateInvoicePDFHandler creates a PDF from an invoice ID and stores the file in firebase.
func GenerateInvoicePDFHandler(c *gin.Context) {
	ctx := context.Background()
	invoiceID := c.Param("id")

	dbApp, ok := c.MustGet("databaseConn").(*db.App)
	if !ok {
		c.AbortWithError(
			http.StatusInternalServerError,
			errors.New("unable to connect to database."))
	}

	req := db.InvoiceRequest{ID: invoiceID}
	invoice, err := dbApp.GetInvoice(ctx, &req)
	if err != nil {
		log.Println(errors.ErrorStack(err))
		c.AbortWithError(http.StatusBadRequest, err)
	}

	filePath := fmt.Sprintf("invoices/%s.pdf", invoiceID)
	pdfConstructor := pdf.PDFConstructor{
		Invoice:    invoice,
		OutputPath: filePath,
	}

	if err := pdf.Construct(pdfConstructor); err != nil {
		log.Println(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	if err := dbApp.UploadFile(ctx, invoiceID, filePath); err != nil {
		log.Println(errors.ErrorStack(err))
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

type EmailInvoiceRequest struct {
	ID string `json:"invoice_id" binding:"required"`
}

func EmailInvoice(c *gin.Context) {
	ctx := context.Background()

	var req EmailInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	userID, ok := c.MustGet("userID").(string)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("could not get user ID"))
	}

	dbApp, ok := c.MustGet("databaseConn").(*db.App)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("could not get database connection"))
	}

	user, err := dbApp.GetUser(ctx, &db.UserRequest{ID: userID})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	invoice, err := dbApp.GetInvoice(ctx, &db.InvoiceRequest{ID: req.ID})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	if err := email.SendInvoice(ctx, user, invoice); err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}
}

type GetInvoiceRequest struct {
	ID string `json:"invoice_id" binding:"required"`
}

func GetInvoice(c *gin.Context) {
	ctx := context.Background()

	var req GetInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	dbApp, ok := c.MustGet("databaseConn").(*db.App)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("could not get database connection"))
	}

	invoice, err := dbApp.GetInvoice(ctx, &db.InvoiceRequest{ID: req.ID})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

type NewInvoiceRequest struct {
	ID string `json:"invoice_id" binding:"required"`
}

func NewInvoice(c *gin.Context) {
	ctx := context.Background()

	var req NewInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	dbApp, ok := c.MustGet("databaseConn").(*db.App)
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, errors.New("could not get database connection"))
	}

	invoice, err := dbApp.GetInvoice(ctx, &db.InvoiceRequest{ID: req.ID})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, errors.Trace(err))
	}

	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}
