package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/email"
	"github.com/rstorr/wham-platform/pdf"
)

func main() {

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.Run()
}

func testEmail() {
	email.Send()
}

func testPdf() {
	ctx := context.Background()
	dbApp, err := db.InitDB(ctx)
	if err != nil {
		log.Fatal(errors.ErrorStack(err))
	}
	defer dbApp.CloseDB()

	req := db.InvoiceRequest{ID: "q7C3BSLSELWe1VObYtXp"}
	invoice, err := dbApp.GetInvoice(ctx, &req)

	pdfConstructor := pdf.PDFConstructor{
		Invoice: invoice,
	}
	if err != nil {
		log.Println(errors.ErrorStack(err))
		return
	}

	pdf.Construct(pdfConstructor)

}
