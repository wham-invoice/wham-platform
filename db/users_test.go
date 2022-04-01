package db_test

import (
	"context"

	"github.com/wham-invoice/wham-platform/tests/setup"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type UsersSuite struct {
	setup.ApplicationSuiteCore
}

var _ = gc.Suite(&UsersSuite{})

func (s *UsersSuite) TestUsersInsertAndGet(c *gc.C) {
	ctx := context.Background()
	u := s.AddUser(ctx, c)

	getInvoice, err := s.App.User(ctx, u.ID)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(getInvoice, jc.DeepEquals, u)
}
