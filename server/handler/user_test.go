package handler_test

import (
	"context"

	"github.com/wham-invoice/wham-platform/db"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type usersSuite struct {
	APISuiteCore

	user *db.User
}

var _ = gc.Suite(&usersSuite{})

func (s *usersSuite) SetUpTest(c *gc.C) {
	ctx := context.Background()
	s.APISuiteCore.SetUpTest(c)

	s.user = s.AddUser(ctx, c)
}

func (s *usersSuite) TestUserSummary(c *gc.C) {

	// Make the request and check the results.
	body := s.Get200(c, "/user/summary")
	c.Check(body, jc.JSONEquals, map[string]interface{}{
		"contacts": []interface{}{
			map[string]interface{}{
				"id": s.user.ID,
			},
		},
	})
}
