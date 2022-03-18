package db_test

import (
	"context"

	"github.com/rstorr/wham-platform/tests/setup"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type ContactsSuite struct {
	setup.ApplicationSuiteCore
}

var _ = gc.Suite(&ContactsSuite{})

func (s *ContactsSuite) TestContactInsertAndGet(c *gc.C) {
	ctx := context.Background()
	contact := s.AddContact(ctx, c)

	getInvoice, err := s.App.Contact(ctx, contact.ID)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(contact, jc.DeepEquals, getInvoice)
}
