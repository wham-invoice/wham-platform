package email

import (
	"github.com/rstorr/wham-platform/db"
)

type InvoiceEmail struct {
	i       *db.Invoice
	pdfPath string
}

func Send() error {
	return nil
}
