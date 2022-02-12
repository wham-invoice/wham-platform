package db

import (
	"context"
	"time"

	"github.com/juju/errors"
)

type Invoice struct {
	ID          string
	User        *User
	Client      *User
	Number      int       `firestore:"number"`
	Rate        float32   `firestore:"rate"`
	Hours       float32   `firestore:"hours"`
	Description string    `firestore:"description"`
	IssueDate   time.Time `firestore:"issue_date"`
	DueDate     time.Time `firestore:"due_date"`
	Paid        bool      `firestore:"paid"`
}

type InvoiceRequest struct {
	ID string
}

const invoiceCollection = "invoices"

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

	dataMap := result.Data()

	userID, err := userIDFromInvoiceResult(dataMap)
	if err != nil {
		return invoice, errors.Trace(err)
	}
	user, err := app.GetUser(ctx, &UserRequest{ID: userID})
	if err != nil {
		return invoice, errors.Trace(err)
	}
	invoice.User = user

	clientID, err := clientIDFromInvoiceResult(dataMap)
	if err != nil {
		return invoice, errors.Trace(err)
	}
	client, err := app.GetContact(ctx, &UserRequest{ID: clientID})
	if err != nil {
		return invoice, errors.Trace(err)
	}
	invoice.Client = client

	return invoice, nil
}

func (i *Invoice) GetSubtotal() float32 {
	return i.Hours * i.Rate
}

const gstRate = 0.15

func (i *Invoice) GetGST() float32 {
	return i.Hours * i.Rate * gstRate
}

func (i *Invoice) GetTotal() float32 {
	return i.GetSubtotal() + i.GetGST()
}

// PutPdf uploads the invoice PDF to the DB.
func PutPdf() error {
	return nil
}

func clientIDFromInvoiceResult(data map[string]interface{}) (string, error) {
	var id string
	var ok bool

	if x, found := data["contact_id"]; found {
		if id, ok = x.(string); !ok {
			return "", errors.New("contact_id was not a string.")
		}
	} else {
		return "", errors.New("contact_id not found.")
	}

	return id, nil
}

func userIDFromInvoiceResult(data map[string]interface{}) (string, error) {
	var id string
	var ok bool

	if x, found := data["user_id"]; found {
		if id, ok = x.(string); !ok {
			return "", errors.New("user_id was not a string.")
		}
	} else {
		return "", errors.New("user_id not found.")
	}

	return id, nil
}
