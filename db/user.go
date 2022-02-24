package db

import (
	"context"
	"fmt"

	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/server/auth"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	ID           string       `json:"id"`
	FirstName    string       `firestore:"first_name" json:"first_name"`
	LastName     string       `firestore:"last_name" json:"last_name"`
	Email        string       `firestore:"email" json:"email"`
	AccessToken  string       `firestore:"access_token" json:"access_token"`
	RefreshToken string       `firestore:"refresh_token" json:"refresh_token"`
	OAuth        oauth2.Token `json:"oauth_token"`
}

type UserRequest struct {
	ID string
}

const usersCollection = "users"

func (app *App) GetUser(ctx context.Context, req *UserRequest) (*User, error) {
	var user = new(User)

	result, err := app.firestoreClient.Collection(usersCollection).Doc(req.ID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil
		}
		return user, errors.Trace(err)
	}

	user.ID = result.Ref.ID
	if err := result.DataTo(&user); err != nil {
		return user, errors.Trace(err)
	}

	fmt.Println(user)

	return user, nil
}

func (app *App) AddUser(ctx context.Context, user *User) error {
	_, err := app.firestoreClient.Collection(usersCollection).Doc(
		user.ID).Set(ctx, user)

	return errors.Trace(err)
}

func (u User) GetFullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (app *App) CreateAndSaveUser(
	ctx context.Context,
	uid string,
	info auth.UserInfo,
	authToken oauth2.Token,
) error {

	user := &User{
		ID:           uid,
		FirstName:    info.GivenName,
		LastName:     info.FamilyName,
		Email:        info.Email,
		AccessToken:  authToken.AccessToken,
		RefreshToken: authToken.RefreshToken,
		OAuth:        authToken,
	}
	return app.AddUser(ctx, user)
}
