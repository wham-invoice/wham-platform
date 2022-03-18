package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/juju/errors"
	"google.golang.org/api/iterator"
)

var ContactNotFound = errors.New("contact not found")

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

		contact.ID = doc.Ref.ID
		if err := doc.DataTo(&contact); err != nil {
			return contacts, errors.Trace(err)
		}

		address, err := addressfromUserDoc(doc)
		if err != nil {
			return contacts, errors.Trace(err)
		}
		contact.Address = address

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
