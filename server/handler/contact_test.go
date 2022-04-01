package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/wham-invoice/wham-platform/db"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type ContactsSuite struct {
	APISuiteCore
}

var _ = gc.Suite(&ContactsSuite{})

func (s *ContactsSuite) SetUpTest(c *gc.C) {
	s.APISuiteCore.SetUpTest(c)
}

func (s *ContactsSuite) TestNewContactThenGet(c *gc.C) {
	firstName := strconv.Itoa(rand.Int())
	lastName := strconv.Itoa(rand.Int())
	phone := strconv.Itoa(rand.Int())
	email := strconv.Itoa(rand.Int())
	company := strconv.Itoa(rand.Int())
	adrLine1 := strconv.Itoa(rand.Int())
	adrLine2 := strconv.Itoa(rand.Int())
	suburb := strconv.Itoa(rand.Int())
	postcode := strconv.Itoa(rand.Int())
	country := strconv.Itoa(rand.Int())

	payload, err := json.Marshal(map[string]interface{}{
		"first_name":          firstName,
		"last_name":           lastName,
		"phone":               phone,
		"email":               email,
		"company":             company,
		"address_first_line":  adrLine1,
		"address_second_line": adrLine2,
		"suburb":              suburb,
		"postcode":            postcode,
		"country":             country,
	})
	c.Assert(err, jc.ErrorIsNil)

	getContact := s.Post200(
		c,
		"/contact/new",
		string(payload),
	)

	unmarsh := db.Contact{}
	json.Unmarshal([]byte(getContact), &unmarsh)

	c.Check(getContact, jc.JSONEquals,
		map[string]interface{}{
			"id":         unmarsh.ID,
			"user_id":    s.user.ID,
			"first_name": firstName,
			"last_name":  lastName,
			"phone":      phone,
			"email":      email,
			"company":    company,
			"address": map[string]interface{}{
				"address_first_line":  adrLine1,
				"address_second_line": adrLine2,
				"address_suburb":      suburb,
				"address_postcode":    postcode,
				"address_country":     country,
			},
		})

}

func (s *ContactsSuite) TestContact(c *gc.C) {
	ctx := context.Background()
	contact := s.AddContact(ctx, c, s.user.ID)

	getContact := s.Get200(
		c,
		fmt.Sprintf("/contact/get/%s", contact.ID),
	)

	c.Check(getContact, jc.JSONEquals,
		map[string]interface{}{
			"id":         contact.ID,
			"user_id":    contact.UserID,
			"first_name": contact.FirstName,
			"last_name":  contact.LastName,
			"phone":      contact.Phone,
			"email":      contact.Email,
			"company":    contact.Company,
			"address": map[string]interface{}{
				"address_first_line":  contact.Address.FirstLine,
				"address_second_line": contact.Address.SecondLine,
				"address_suburb":      contact.Address.Suburb,
				"address_postcode":    contact.Address.Postcode,
				"address_country":     contact.Address.Country,
			},
		})
}

func (s *ContactsSuite) TestContacts(c *gc.C) {
	var contacts []db.Contact

	ctx := context.Background()
	contact1 := s.AddContact(ctx, c, s.user.ID)
	contacts = append(contacts, *contact1)
	contact2 := s.AddContact(ctx, c, s.user.ID)
	contacts = append(contacts, *contact2)

	getContacts := s.Get200(c, "/user/contacts")

	contactsRespList := []db.Contact{}
	json.Unmarshal([]byte(getContacts), &contactsRespList)

	c.Assert(contactsRespList, gc.HasLen, 2)

	for _, contact := range contacts {
		c.Check(contactsRespList, jc.Contains, contact)
	}
}
