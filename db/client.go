package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/juju/errors"
)

type User struct {
	ID        string
	FirstName string `firestore:"first_name"`
	LastName  string `firestore:"last_name"`
	Phone     string `firestore:"phone_number"`
	Email     string `firestore:"email"`
	Company   string `firestore:"company"`
	Address   *Address
}

type Address struct {
	FirstLine  string `firestore:"address_first_line"`
	SecondLine string `firestore:"address_second_line"`
	Suburb     string `firestore:"address_suburb"`
	Postcode   string `firestore:"address_postcode"`
	Country    string `firestore:"address_country"`
}

type UserRequest struct {
	ID string
}

const usersCollection = "users"

func (app *App) GetUser(ctx context.Context, req *UserRequest) (*User, error) {
	var user = new(User)

	result, err := app.client.Collection(usersCollection).Doc(req.ID).Get(ctx)
	if err != nil {
		return user, errors.Trace(err)
	}

	user.ID = result.Ref.ID
	if err := result.DataTo(&user); err != nil {
		return user, errors.Trace(err)
	}

	address, err := addressfromUserDoc(result)
	if err != nil {
		return user, errors.Trace(err)
	}
	user.Address = address

	return user, nil
}

func (u User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func addressfromUserDoc(doc *firestore.DocumentSnapshot) (*Address, error) {
	var address = new(Address)
	if err := doc.DataTo(&address); err != nil {
		return address, errors.Trace(err)
	}

	return address, nil
}
