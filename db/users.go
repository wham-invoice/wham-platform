package db

import (
	"context"
	"fmt"

	"github.com/juju/errors"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var UserNotFound = errors.New("user not found")

type User struct {
	ID        string `json:"id"`
	FirstName string `firestore:"first_name" json:"first_name"`
	LastName  string `firestore:"last_name" json:"last_name"`
	Email     string `firestore:"email" json:"email"`
	// TODO we need to keep an eye on oauth2 expiries and refresh tokens when necessary.
	OAuth oauth2.Token `json:"oauth_token"`
}

type UserSummary struct {
	InvoiceTotal float32 `json:"invoice_total"`
	InvoicePaid  float32 `json:"invoice_paid"`
}

const usersCollection = "users"

func (app *App) AddUser(ctx context.Context, user *User) error {
	_, err := app.firestoreClient.Collection(usersCollection).Doc(
		user.ID).Set(ctx, user)

	return errors.Trace(err)
}

func (app *App) User(ctx context.Context, id string) (*User, error) {
	var user = new(User)

	result, err := app.firestoreClient.Collection(usersCollection).Doc(
		id).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return user, errors.Annotatef(err, "errored getting user %s", id)
	}

	user.ID = result.Ref.ID
	if err := result.DataTo(&user); err != nil {
		return user, errors.Trace(err)
	}

	return user, nil
}

func (u User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u User) Invoices(ctx context.Context, app *App) ([]Invoice, error) {
	return app.invoicesForUser(ctx, u.ID)
}

func (u User) Contacts(ctx context.Context, app *App) ([]Contact, error) {
	return app.contactsForUser(ctx, u.ID)
}

func (u User) Summary(ctx context.Context, app *App) (UserSummary, error) {
	var summary UserSummary
	total, paid, err := app.invoiceTotalsForUser(ctx, u.ID)
	if err != nil {
		return summary, errors.Trace(err)
	}

	summary.InvoiceTotal = total
	summary.InvoicePaid = paid

	return summary, nil
}

// Santize returns a copy of the user with sensitive info removed.
func (u User) Sanitize() User {
	u.OAuth = oauth2.Token{}

	return u
}

type UserInfo struct {
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Name       string `json:"name"`
}

func (app *App) NewUser(
	ctx context.Context,
	uid string,
	info UserInfo,
	authToken oauth2.Token,
) error {

	user := &User{
		ID:        uid,
		FirstName: info.GivenName,
		LastName:  info.FamilyName,
		Email:     info.Email,
		OAuth:     authToken,
	}
	return app.AddUser(ctx, user)
}

func (app *App) UsersDeleteAll(ctx context.Context, batchSize int) error {
	for {
		iter := app.firestoreClient.Collection(usersCollection).Limit(batchSize).Documents(ctx)
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
