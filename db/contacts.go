package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/juju/errors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ContactNotFound = errors.New("contact not found")

type Contact struct {
	ID        string   `firestore:"id" json:"id"`
	UserID    string   `firestore:"user_id" json:"user_id"`
	FirstName string   `firestore:"first_name" json:"first_name"`
	LastName  string   `firestore:"last_name" json:"last_name"`
	Phone     string   `firestore:"phone" json:"phone"`
	Email     string   `firestore:"email" json:"email"`
	Company   string   `firestore:"company" json:"company"`
	Address   *Address `firestore:"address" json:"address"`
}

type Address struct {
	FirstLine  string `firestore:"address_first_line" json:"address_first_line"`
	SecondLine string `firestore:"address_second_line" json:"address_second_line"`
	Suburb     string `firestore:"address_suburb" json:"address_suburb"`
	Postcode   string `firestore:"address_postcode" json:"address_postcode"`
	Country    string `firestore:"address_country" json:"address_country"`
}

const contactsCollection = "contacts"

func (app *App) AddContact(ctx context.Context, contact *Contact) (string, error) {
	ref, _, err := app.firestoreClient.Collection(contactsCollection).Add(ctx, contact)
	if err != nil {
		return "", errors.Trace(err)
	}

	id := ref.ID

	return id, nil
}

func (app *App) Contact(ctx context.Context, id string) (*Contact, error) {
	var contact = new(Contact)

	result, err := app.firestoreClient.Collection(contactsCollection).Doc(id).Get(ctx)
	if err != nil {
		return contact, errors.Trace(err)
	}

	if !result.Exists() {
		return contact, errors.Errorf("contact with ID does not exist: %s" + id)
	}

	if err := result.DataTo(&contact); err != nil {
		return contact, errors.Trace(err)
	}

	contact.ID = result.Ref.ID

	return contact, nil
}

// TODO: test
func (c *Contact) Delete(ctx context.Context, app *App) error {
	_, err := app.firestoreClient.Collection(contactsCollection).Doc(c.ID).Delete(ctx)
	if status.Code(err) == codes.NotFound {
		return ContactNotFound
	}

	return errors.Trace(err)
}

func (app *App) contactsForUser(ctx context.Context, userID string) ([]Contact, error) {
	contacts := []Contact{}

	iter := app.firestoreClient.Collection(contactsCollection).Where("user_id", "==", userID).Documents(ctx)
	for {
		var contact = new(Contact)
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return contacts, err
		}

		if err := doc.DataTo(&contact); err != nil {
			return contacts, errors.Trace(err)
		}

		contact.ID = doc.Ref.ID

		contacts = append(contacts, *contact)

	}

	return contacts, nil
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

func (app *App) ContactsDeleteAll(ctx context.Context, batchSize int) error {
	for {
		iter := app.firestoreClient.Collection(contactsCollection).Limit(batchSize).Documents(ctx)
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
