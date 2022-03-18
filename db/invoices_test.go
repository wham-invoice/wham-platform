package db_test

import (
	"context"

	"github.com/rstorr/wham-platform/tests/setup"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type InvoicesSuite struct {
	setup.ApplicationSuiteCore
}

var _ = gc.Suite(&InvoicesSuite{})

func (s *InvoicesSuite) TestInvoiceInsertAndGet(c *gc.C) {
	ctx := context.Background()
	inv := s.AddInvoice(ctx, c)

	getInvoice, err := s.App.Invoice(ctx, inv.ID)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(inv, jc.DeepEquals, getInvoice)
}
