package db

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/juju/errors"
	"google.golang.org/api/option"
)

type App struct {
	client firestore.Client
}

func InitDB(ctx context.Context) (*App, error) {
	var app = new(App)

	opt := option.WithCredentialsFile("/Users/work/go/src/github.com/rstorr/wham-platform/secrets/firebase_service_account_key.json")
	firebaseApp, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return app, errors.Trace(err)
	}
	client, err := firebaseApp.Firestore(ctx)
	if err != nil {
		return app, errors.Trace(err)
	}
	app.client = *client

	return app, nil

}

func (a *App) CloseDB() {
	a.client.Close()
}
