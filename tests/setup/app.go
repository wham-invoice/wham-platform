package setup

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/wham-invoice/wham-platform/db"
	"github.com/wham-invoice/wham-platform/util"

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

// TODO we're deleting all tests here!
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
	id := strconv.Itoa(rand.Int())
	firstName := strconv.Itoa(rand.Int())
	lastName := strconv.Itoa(rand.Int())
	email := strconv.Itoa(rand.Int())
	return &db.User{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		OAuth:     oauth2.Token{},
	}
}

func (s *ApplicationSuiteCore) AddInvoice(
	c *gc.C,
	userID string,
) *db.Invoice {
	ctx := context.Background()

	invoice := CreateInvoice(userID)
	id, err := s.App.AddInvoice(ctx, invoice)
	c.Assert(err, jc.ErrorIsNil)
	invoice.ID = id

	return invoice
}

func CreateInvoice(userID string) *db.Invoice {
	contactID := strconv.Itoa(rand.Int())
	pdfID := strconv.Itoa(rand.Int())
	number := rand.Int()
	hours := rand.Float32()
	rate := rand.Float32()
	description := strconv.Itoa(rand.Int())
	issueDate := time.Now()
	dueDate := time.Now().Add(time.Hour * time.Duration(240))

	return &db.Invoice{
		UserID:      userID,
		ContactID:   contactID,
		PDFID:       pdfID,
		Number:      number,
		Rate:        rate,
		Hours:       hours,
		Description: description,
		IssueDate:   issueDate,
		DueDate:     dueDate,
	}
}

func (s *ApplicationSuiteCore) AddContact(
	ctx context.Context,
	c *gc.C,
	userID string,
) *db.Contact {
	contact := CreateContact(userID)
	id, err := s.App.AddContact(ctx, &contact)
	c.Assert(err, jc.ErrorIsNil)
	contact.ID = id

	return &contact
}

func CreateContact(userID string) db.Contact {
	firstName := strconv.Itoa(rand.Int())
	lastName := strconv.Itoa(rand.Int())
	email := strconv.Itoa(rand.Int())
	phone := strconv.Itoa(rand.Int())
	company := strconv.Itoa(rand.Int())
	address := &db.Address{
		FirstLine:  strconv.Itoa(rand.Int()),
		SecondLine: strconv.Itoa(rand.Int()),
		Suburb:     strconv.Itoa(rand.Int()),
		Postcode:   strconv.Itoa(rand.Int()),
		Country:    strconv.Itoa(rand.Int()),
	}

	return db.Contact{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Company:   company,
		Address:   address,
	}
}
