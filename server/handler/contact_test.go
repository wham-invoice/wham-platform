package handler_test

import (
	"context"
	"fmt"

	"github.com/rstorr/wham-platform/db"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type contactsSuite struct {
	APISuiteCore

	contact *db.Contact
}

var _ = gc.Suite(&contactsSuite{})

func (s *contactsSuite) SetUpTest(c *gc.C) {
	ctx := context.Background()
	s.APISuiteCore.SetUpTest(c)

	s.contact = s.AddContact(ctx, c)
}

func (s *contactsSuite) TestContacts(c *gc.C) {

	// Make the request and check the results.
	body := s.Get200(c, "/user/contacts")
	c.Check(body, jc.JSONEquals, map[string]interface{}{
		"contacts": []interface{}{
			map[string]interface{}{
				"id": s.contact.ID,
			},
		},
	})
}

func (s *contactsSuite) TestContact(c *gc.C) {

	// Make the request and check the results.
	body := s.Get200(c, fmt.Sprintf("/contact/get/", s.contact.ID))
	c.Check(body, jc.JSONEquals, map[string]interface{}{
		"contacts": []interface{}{
			map[string]interface{}{
				"id": s.contact.ID,
			},
		},
	})
}
