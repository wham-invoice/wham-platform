package main

import (
	"context"
	"log"

	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/pdf"
)

func main() {

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
