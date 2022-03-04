package db

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/juju/errors"
)

type Invoice struct {
	ID          string
	UserID      string    `firestore:"user_id"`
	ContactID   string    `firestore:"contact_id"`
	PDFID       string    `firestore:"pdf_id"`
	Number      int       `firestore:"number"`
	Rate        float32   `firestore:"rate"`
	Hours       float32   `firestore:"hours"`
	Description string    `firestore:"description"`
	IssueDate   time.Time `firestore:"issue_date"`
	DueDate     time.Time `firestore:"due_date"`
	Paid        bool      `firestore:"paid"`
	URLCode     string    `firestore:"url_code"`
}

type InvoiceRequest struct {
	ID string
}

const invoiceCollection = "invoices"

func (app *App) AddInvoice(ctx context.Context, invoice *Invoice) error {

	_, _, err := app.firestoreClient.Collection(invoiceCollection).Add(ctx, invoice)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (app *App) GetInvoice(ctx context.Context, req *InvoiceRequest) (*Invoice, error) {
	var invoice = new(Invoice)

	result, err := app.firestoreClient.Collection(invoiceCollection).Doc(req.ID).Get(ctx)
	if err != nil {
		return invoice, errors.Trace(err)
	}

	invoice.ID = result.Ref.ID

	if err := result.DataTo(&invoice); err != nil {
		return invoice, errors.Trace(err)
	}

	return invoice, nil
}

func (i *Invoice) GetUser(ctx context.Context, app *App) (*User, error) {
	return app.GetUser(ctx, &UserRequest{ID: i.UserID})
}

func (i *Invoice) GetContact(ctx context.Context, app *App) (*Contact, error) {
	return app.GetContact(ctx, &ContactRequest{ID: i.ContactID})
}

func (i *Invoice) GetSubtotal() float32 {
	return i.Hours * i.Rate
}

func (i *Invoice) GetGST() float32 {
	const gstRate = 0.15
	return i.Hours * i.Rate * gstRate
}

func (i *Invoice) GetTotal() float32 {
	return i.GetSubtotal() + i.GetGST()
}

func (app *App) UploadPDFDeleteLocal(ctx context.Context, fileName, filePath string) error {

	bucket, err := app.storageClient.Bucket("invoice_pdf")
	if err != nil {
		return errors.Trace(err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	writer := bucket.Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(writer, f); err != nil {
		return errors.Trace(err)
	}
	if err := writer.Close(); err != nil {
		return errors.Trace(err)
	}

	return nil
}
