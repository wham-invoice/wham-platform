package db_test

import (
	"context"

	"github.com/wham-invoice/wham-platform/db"
	"github.com/wham-invoice/wham-platform/tests/setup"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type ContactsSuite struct {
	setup.ApplicationSuiteCore

	user *db.User
}

var _ = gc.Suite(&ContactsSuite{})

func (s *ContactsSuite) SetUpTest(c *gc.C) {

	s.user = s.AddUser(context.Background(), c)
}

func (s *ContactsSuite) TestContactInsertAndGet(c *gc.C) {
	ctx := context.Background()
	contact := s.AddContact(ctx, c, s.user.ID)

	getContact, err := s.App.Contact(ctx, contact.ID)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(getContact, jc.DeepEquals, contact)
}
