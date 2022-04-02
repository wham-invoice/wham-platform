package db

import (
	"context"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"github.com/juju/errors"
	"google.golang.org/api/option"
)

type App struct {
	firestoreClient *firestore.Client
	storageClient   *storage.Client
}

// TODO DB config should be stored in config file.
func Init(ctx context.Context) (*App, error) {
	var app = new(App)

	config := &firebase.Config{
		StorageBucket: "wham-ad61b.appspot.com",
	}
	opt := option.WithCredentialsFile("./secrets/firebase_service_account_key.json")
	firebaseApp, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, errors.Trace(err)
	}

	fs, err := firebaseApp.Firestore(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	app.firestoreClient = fs

	storage, err := firebaseApp.Storage(ctx)
	if err != nil {
		return nil, errors.Trace(err)
	}
	app.storageClient = storage

	return app, nil
}

func (a *App) CloseDB() {
	a.firestoreClient.Close()
}
