package db

import (
	"context"
	"fmt"

	"github.com/juju/errors"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	ID        string       `json:"id"`
	FirstName string       `firestore:"first_name" json:"first_name"`
	LastName  string       `firestore:"last_name" json:"last_name"`
	Email     string       `firestore:"email" json:"email"`
	OAuth     oauth2.Token `json:"oauth_token"`
}

const usersCollection = "users"

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

func (app *App) AddUser(ctx context.Context, user *User) error {
	_, err := app.firestoreClient.Collection(usersCollection).Doc(
		user.ID).Set(ctx, user)

	return errors.Trace(err)
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

type UserInfo struct {
	Email      string `json:"email"`
	FamilyName string `json:"family_name"`
	GivenName  string `json:"given_name"`
	Name       string `json:"name"`
	// Sub        string `json: "sub"`
	// An identifier for the user,
	// unique among all Google accounts and never reused.
	// A Google account can have multiple email addresses at different points in time, but the sub value is never changed.
	// Use sub within your application as the unique-identifier key for the user
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
