package db

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/juju/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var InvoiceNotFound = errors.New("invoice not found")

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

	doc, err := app.firestoreClient.Collection(invoiceCollection).Doc(req.ID).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return invoice, InvoiceNotFound
	}

	if err != nil {
		return invoice, errors.Trace(err)
	}

	invoice.ID = doc.Ref.ID

	if err := doc.DataTo(&invoice); err != nil {
		return invoice, errors.Trace(err)
	}

	return invoice, nil
}

func (app *App) GetInvoicesForUser(ctx context.Context, userID string) ([]Invoice, error) {
	invoices := []Invoice{}

	iter := app.firestoreClient.Collection(invoiceCollection).Where("user_id", "==", userID).Documents(ctx)
	for {
		var invoice = new(Invoice)
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return invoices, err
		}

		invoice.ID = doc.Ref.ID

		if err := doc.DataTo(&invoice); err != nil {
			return invoices, errors.Trace(err)
		}

		invoices = append(invoices, *invoice)
	}

	return invoices, nil
}

func (i *Invoice) GetUser(ctx context.Context, app *App) (*User, error) {
	return app.GetUser(ctx, i.UserID)
}

func (i *Invoice) GetContact(ctx context.Context, app *App) (*Contact, error) {
	return app.GetContact(ctx, i.ContactID)
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
