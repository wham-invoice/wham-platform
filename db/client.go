package db

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/juju/errors"
)

type Client struct {
	ID        string
	FirstName string `firestore:"first_name"`
	LastName  string `firestore:"last_name"`
	Email     string `firestore:"email"`
	Address   *Address
}

type Address struct {
	FirstLine  string `firestore:"second_line"`
	SecondLine string `firestore:"first_line"`
	Suburb     string `firestore:"suburb"`
	Postcode   string `firestore:"postcode"`
	Country    string `firestore:"country"`
}

type ClientRequest struct {
	ID string
}

const clientCollection = "clients"

func (app *App) GetClient(ctx context.Context, req *ClientRequest) (*Client, error) {
	var client = new(Client)

	result, err := app.client.Collection(clientCollection).Doc(req.ID).Get(ctx)
	if err != nil {
		return client, errors.Trace(err)
	}

	client.ID = result.Ref.ID
	fmt.Printf("Client Name1: %v \n", result.Data()["first_name"])
	if err := result.DataTo(&client); err != nil {
		return client, errors.Trace(err)
	}

	address, err := addressfromClientDoc(result)
	if err != nil {
		return client, errors.Trace(err)
	}
	client.Address = address

	fmt.Printf("Client Name: %v \n", result.Data()["first_name"])
	fmt.Printf("Client: %v\n", client)
	fmt.Printf("Client obj name: %v\n", client.FirstName)
	return client, nil
}

func addressfromClientDoc(doc *firestore.DocumentSnapshot) (*Address, error) {
	var address = new(Address)
	if err := doc.DataTo(&address); err != nil {
		return address, errors.Trace(err)
	}

	return address, nil
}
