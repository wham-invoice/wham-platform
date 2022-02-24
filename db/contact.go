package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/juju/errors"
)

type Contact struct {
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

type ContactRequest struct {
	ID string
}

const contactsCollection = "contacts"

func (app *App) GetContact(ctx context.Context, req *ContactRequest) (*Contact, error) {
	var contact = new(Contact)

	result, err := app.firestoreClient.Collection(contactsCollection).Doc(req.ID).Get(ctx)
	if err != nil {
		return contact, errors.Trace(err)
	}

	contact.ID = result.Ref.ID
	if err := result.DataTo(&contact); err != nil {
		return contact, errors.Trace(err)
	}

	address, err := addressfromUserDoc(result)
	if err != nil {
		return contact, errors.Trace(err)
	}
	contact.Address = address

	return contact, nil
}

func addressfromUserDoc(doc *firestore.DocumentSnapshot) (*Address, error) {
	var address = new(Address)
	if err := doc.DataTo(&address); err != nil {
		return address, errors.Trace(err)
	}

	return address, nil
}

func (c Contact) GetFullName() string {
	return fmt.Sprintf("%s %s", c.FirstName, c.LastName)
}

