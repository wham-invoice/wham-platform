package db_test

import (
	"context"

	"github.com/wham-invoice/wham-platform/db"
	"github.com/wham-invoice/wham-platform/tests/setup"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type InvoicesSuite struct {
	setup.ApplicationSuiteCore

	user *db.User
}

var _ = gc.Suite(&InvoicesSuite{})

func (s *InvoicesSuite) SetUpTest(c *gc.C) {

	s.user = s.AddUser(context.Background(), c)
}

func (s *InvoicesSuite) TestInvoiceInsertAndGet(c *gc.C) {
	inv := s.AddInvoice(c, s.user.ID)

	getInvoice, err := s.App.Invoice(context.Background(), inv.ID)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(getInvoice, jc.DeepEquals, inv)
}
