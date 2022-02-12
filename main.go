package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/pdf"
)

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(dbApp *db.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", dbApp)
		c.Next()
	}
}

func main() {
	ctx := context.Background()

	dbApp, err := db.Init(ctx)
	if err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
	defer dbApp.CloseDB()

	router := gin.Default()
	router.Use(ApiMiddleware(dbApp))
	router.GET("/invoice/:id/pdf", GenerateInvoicePDFHandler)
	router.Run(":8080")
}

// GenerateInvoicePDFHandler creates a PDF from an invoice ID and stores the file in firebase.
func GenerateInvoicePDFHandler(c *gin.Context) {
	invoiceID := c.Param("id")

	dbApp, ok := c.MustGet("databaseConn").(*db.App)
	if !ok {
		// handle error here...
	}

	ctx := context.Background()

	req := db.InvoiceRequest{ID: invoiceID}
	invoice, err := dbApp.GetInvoice(ctx, &req)

	filePath := "test_dir/thefile.pdf"
	pdfConstructor := pdf.PDFConstructor{
		Invoice:    invoice,
		OutputPath: filePath,
	}
	if err != nil {
		log.Println(errors.ErrorStack(err))
	}

	if err := pdf.Construct(pdfConstructor); err != nil {
		log.Println(errors.ErrorStack(err))
	}

	if err := dbApp.UploadFile(ctx, "yozaFile", filePath); err != nil {
		log.Println(errors.ErrorStack(err))
	}

}
