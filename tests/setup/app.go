package setup

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/util"

	jc "github.com/juju/testing/checkers"
	"golang.org/x/oauth2"
	gc "gopkg.in/check.v1"
)

type ApplicationSuiteCore struct {
	App *db.App
}

func (s *ApplicationSuiteCore) SetUpSuite(c *gc.C) {
	ctx := context.Background()
	c.Assert(util.SetDebugLogger(), jc.ErrorIsNil)
	app, err := db.Init(ctx)
	c.Assert(err, jc.ErrorIsNil)
	s.App = app
}

func (s *ApplicationSuiteCore) SetUpTest(c *gc.C) {
	ctx := context.Background()
	c.Assert(util.SetDebugLogger(), jc.ErrorIsNil)
	c.Assert(s.App.UsersDeleteAll(ctx, 50), jc.ErrorIsNil)
	c.Assert(s.App.InvoicesDeleteAll(ctx, 50), jc.ErrorIsNil)
	c.Assert(s.App.ContactsDeleteAll(ctx, 50), jc.ErrorIsNil)
	// TODO delete all files from storage.
}

func (s *ApplicationSuiteCore) TearDownSuite(c *gc.C) {
	// TODO close firestore db conn? delete firestore test db?
}

type UserFunc func(*db.User)

func (s *ApplicationSuiteCore) AddUser(
	ctx context.Context,
	c *gc.C,
) *db.User {
	user := CreateUser()
	err := s.App.AddUser(ctx, user)
	c.Assert(err, jc.ErrorIsNil)
	return user
}

func CreateUser() *db.User {
	firstName := strconv.Itoa(rand.Int())
	lastName := strconv.Itoa(rand.Int())
	email := strconv.Itoa(rand.Int())
	return &db.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		OAuth:     oauth2.Token{},
	}
}

func (s *ApplicationSuiteCore) AddInvoice(
	ctx context.Context,
	c *gc.C,
) *db.Invoice {
	invoice := CreateInvoice()
	id, err := s.App.AddInvoice(ctx, invoice)
	c.Assert(err, jc.ErrorIsNil)
	invoice.ID = id

	return invoice
}

func CreateInvoice() *db.Invoice {
	firstName := strconv.Itoa(rand.Int())
	lastName := strconv.Itoa(rand.Int())
	email := strconv.Itoa(rand.Int())
	return &db.Invoice{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		OAuth:     oauth2.Token{},
	}
}

func (s *ApplicationSuiteCore) AddContact(
	ctx context.Context,
	c *gc.C,
) *db.Contact {
	contact := CreateContact()
	id, err := s.App.AddContact(ctx, contact)
	c.Assert(err, jc.ErrorIsNil)
	contact.ID = id

	return contact
}

func CreateContact() *db.Contact {
	firstName := strconv.Itoa(rand.Int())
	lastName := strconv.Itoa(rand.Int())
	email := strconv.Itoa(rand.Int())
	return &db.Contact{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		OAuth:     oauth2.Token{},
	}
}
