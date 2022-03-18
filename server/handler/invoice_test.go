package handler_test

import (
	"context"
	"fmt"

	"github.com/rstorr/wham-platform/db"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type invoicesSuite struct {
	APISuiteCore

	invoice *db.Invoice
}

var _ = gc.Suite(&invoicesSuite{})

func (s *invoicesSuite) SetUpTest(c *gc.C) {
	ctx := context.Background()
	s.APISuiteCore.SetUpTest(c)

	s.invoice = s.AddInvoice(ctx, c)
}

func (s *invoicesSuite) TestInvoices(c *gc.C) {

	body := s.Get200(c, "/user/invoices")
	c.Check(body, jc.JSONEquals, map[string]interface{}{
		"contacts": []interface{}{
			map[string]interface{}{
				"id": s.invoice.ID,
			},
		},
	})
}

func (s *invoicesSuite) TestInvoice(c *gc.C) {

	body := s.Get200(c, fmt.Sprintf("/invoice/get/", s.invoice.ID))
	c.Check(body, jc.JSONEquals, map[string]interface{}{
		"contacts": []interface{}{
			map[string]interface{}{
				"id": s.invoice.ID,
			},
		},
	})
}

func (s *invoicesSuite) TestInvoiceEmail(c *gc.C) {
	s.Post200(c, "/invoice/email")
}

func (s *invoicesSuite) TestInvoiceEmail400(c *gc.C) {
	s.Post400(c, "/invoice/email")
}

func (s *invoicesSuite) TestInvoiceNew(c *gc.C) {
	s.Post200(c, "/invoice/new")
}

func (s *invoicesSuite) TestInvoiceNew400(c *gc.C) {
	s.Post400(c, "/invoice/new")
}

func (s *invoicesSuite) TestViewInvoice(c *gc.C) {
	s.Get200(c, fmt.Sprintf("/invoice/view/:%s", s.invoice.ID))
}

func (s *invoicesSuite) TestUserInvoice(c *gc.C) {
	s.Get200(c, "/user/invoices")
}
