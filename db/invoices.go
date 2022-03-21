package db

import (
	"context"
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

type InvoiceDetail struct {
	PDFID       string
	User        *User
	Contact     *Contact
	Number      int
	Rate        float32
	Hours       float32
	Description string
	IssueDate   time.Time
	DueDate     time.Time
	Paid        bool
}

const invoicesCollection = "invoices"

func (app *App) AddInvoice(ctx context.Context, invoice *Invoice) (string, error) {
	ref, _, err := app.firestoreClient.Collection(invoicesCollection).Add(ctx, invoice)
	if err != nil {
		return "", errors.Trace(err)
	}

	id := ref.ID

	return id, nil
}

func (app *App) Invoice(ctx context.Context, id string) (*Invoice, error) {
	var invoice = new(Invoice)

	doc, err := app.firestoreClient.Collection(invoicesCollection).Doc(id).Get(ctx)
	if status.Code(err) == codes.NotFound {
		return invoice, InvoiceNotFound
	}

	if err != nil {
		return invoice, errors.Trace(err)
	}

	if err := doc.DataTo(&invoice); err != nil {
		return invoice, errors.Trace(err)
	}

	invoice.ID = doc.Ref.ID

	return invoice, nil
}

func (i *Invoice) Detail(ctx context.Context, app *App) (*InvoiceDetail, error) {
	user, err := i.User(ctx, app)
	if err != nil {
		return nil, errors.Trace(err)
	}
	// user without oauth token.
	userSafe := &User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	contact, err := i.Contact(ctx, app)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return &InvoiceDetail{
		PDFID:       i.PDFID,
		User:        userSafe,
		Contact:     contact,
		Number:      i.Number,
		Rate:        i.Rate,
		Hours:       i.Hours,
		Description: i.Description,
		IssueDate:   i.IssueDate,
		DueDate:     i.DueDate,
		Paid:        i.Paid,
	}, nil
}

func (app *App) invoicesForUser(ctx context.Context, userID string) ([]Invoice, error) {
	invoices := []Invoice{}

	iter := app.firestoreClient.Collection(invoicesCollection).Where("user_id", "==", userID).Documents(ctx)
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

func (app *App) invoiceTotalsForUser(ctx context.Context, userID string) (float32, float32, error) {
	var total, paid float32

	invoices, err := app.invoicesForUser(ctx, userID)
	if err != nil {
		return total, paid, errors.Trace(err)
	}

	for _, invoice := range invoices {
		total += invoice.GetTotal()
		if invoice.Paid {
			paid += invoice.GetTotal()
		}
	}

	return total, paid, nil
}

func (i *Invoice) User(ctx context.Context, app *App) (*User, error) {
	return app.User(ctx, i.UserID)
}

func (i *Invoice) Contact(ctx context.Context, app *App) (*Contact, error) {
	return app.Contact(ctx, i.ContactID)
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

func (app *App) InvoicesDeleteAll(ctx context.Context, batchSize int) error {
	for {
		iter := app.firestoreClient.Collection(invoicesCollection).Limit(batchSize).Documents(ctx)
		numDeleted := 0

		batch := app.firestoreClient.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
}
