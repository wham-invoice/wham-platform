package db

import (
	"context"
	"fmt"
	"time"

	"github.com/juju/errors"
)

type Invoice struct {
	ID          string
	Client      *Client
	Rate        float32   `firestore:"rate"`
	Hours       float32   `firestore:"hours"`
	Description string    `firestore:"description"`
	DueDate     time.Time `firestore:"due_date"`
	Paid        bool      `firestore:"paid"`
}

type InvoiceRequest struct {
	ID string
}

const invoiceCollection = "invoices"

func (app *App) GetInvoice(ctx context.Context, req *InvoiceRequest) (*Invoice, error) {
	var invoice = new(Invoice)

	result, err := app.client.Collection(invoiceCollection).Doc(req.ID).Get(ctx)
	if err != nil {
		return invoice, errors.Trace(err)
	}

	invoice.ID = result.Ref.ID
	if err := result.DataTo(&invoice); err != nil {
		return invoice, errors.Trace(err)
	}

	dataMap := result.Data()
	clientID, err := clientIDFromInvoiceResult(dataMap)
	if err != nil {
		return invoice, errors.Trace(err)
	}
	client, err := app.GetClient(ctx, &ClientRequest{ID: clientID})
	if err != nil {
		return invoice, errors.Trace(err)
	}
	invoice.Client = client

	fmt.Printf("invoice: %v\n", invoice)
	return invoice, nil
}

func clientIDFromInvoiceResult(data map[string]interface{}) (string, error) {
	var id string
	var ok bool

	if x, found := data["client_id"]; found {
		if id, ok = x.(string); !ok {
			return "", errors.New("client_id was not a string.")
		}
	} else {
		return "", errors.New("client_id not found.")
	}

	return id, nil
}
