package handler_test

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rstorr/wham-platform/db"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type invoicesSuite struct {
	APISuiteCore

	user *db.User
}

var _ = gc.Suite(&invoicesSuite{})

func (s *invoicesSuite) SetUpTest(c *gc.C) {
	s.APISuiteCore.SetUpTest(c)

	s.user = s.AddUser(context.Background(), c)
}

func (s *invoicesSuite) TestInvoices(c *gc.C) {
	invoice1 := s.AddInvoice(c, s.user.ID)
	_ = s.AddInvoice(c, s.user.ID)

	body := s.Get200(c, "/user/invoices")
	c.Check(body, jc.JSONEquals, map[string]interface{}{
		"invoices": []interface{}{
			map[string]interface{}{
				"id": invoice1.ID,
			},
		},
	})
}

func (s *invoicesSuite) TestInvoice(c *gc.C) {
	invoice := s.AddInvoice(c, s.user.ID)
	body := s.Get200(c, fmt.Sprintf("/invoice/get/%s", invoice.ID))
	c.Check(body, jc.JSONEquals, map[string]interface{}{
		"invoices": []interface{}{
			map[string]interface{}{
				"id": invoice.ID,
			},
		},
	})
}

func (s *invoicesSuite) TestInvoiceEmail(c *gc.C) {
	invoice := s.AddInvoice(c, s.user.ID)

	payload, err := json.Marshal(map[string]interface{}{
		"invoice_id": invoice.ID,
	})
	c.Assert(err, jc.ErrorIsNil)

	s.Post204(c, "/invoice/email", string(payload))
}

func (s *invoicesSuite) TestInvoiceEmail400(c *gc.C) {
	invoice := s.AddInvoice(c, s.user.ID)

	payload, err := json.Marshal(map[string]interface{}{
		"invoice_id": invoice.ID,
	})
	c.Assert(err, jc.ErrorIsNil)

	// TODO 400 ?
	s.Post400(c, "/invoice/email", string(payload))
}

// func (s *invoicesSuite) TestInvoiceNew(c *gc.C) {
// 	s.Post200(c, "/invoice/new")
// }

// func (s *invoicesSuite) TestInvoiceNew400(c *gc.C) {
// 	s.Post400(c, "/invoice/new")
// }

// func (s *invoicesSuite) TestViewInvoice(c *gc.C) {
// 	s.Get200(c, fmt.Sprintf("/invoice/view/:%s", s.invoice.ID))
// }

// func (s *invoicesSuite) TestUserInvoice(c *gc.C) {
// 	s.Get200(c, "/user/invoices")
// }
