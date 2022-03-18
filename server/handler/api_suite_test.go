package handler_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rstorr/wham-platform/db"
	"github.com/rstorr/wham-platform/server/handler"
	"github.com/rstorr/wham-platform/tests/setup"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
)

type APISuiteCore struct {
	setup.ApplicationSuiteCore
	ngin *gin.Engine

	user *db.User
}

func (s *APISuiteCore) SetUpSuite(c *gc.C) {
	s.ApplicationSuiteCore.SetUpSuite(c)

	root, err := handler.Root(handler.Config{
		AllowOrigin: "http://test.origin",
		AppDB:       s.App,
	})
	c.Assert(err, jc.ErrorIsNil)

	ngin := gin.New()
	root.Install(&ngin.RouterGroup)
	s.ngin = ngin
}

func (s *APISuiteCore) SetUpTest(c *gc.C) {
	s.ApplicationSuiteCore.SetUpTest(c)

	s.user = s.AddUser(context.Background(), c)
}

// EnsureDBUser is part of the handler.Github interface.
func (s *APISuiteCore) EnsureDBUser(
	ctx context.Context,
	app *db.App,
) (*db.User, error) {
	if s.user != nil {
		return s.user, nil
	}
	panic("oh no EnsureDBUser: s.user is nil")
}

// RequestOAuthToken is part of the handler.Github interface.
func (s *APISuiteCore) RequestOAuthToken(
	tempCode,
	state string,
) (string, error) {
	// RequestOAuthToken sends tempcode and state to GH and receives the users oauth token.
	return "itstokentime", nil
}

func (s *APISuiteCore) Serve(req *http.Request) *http.Response {
	w := httptest.NewRecorder()
	s.ngin.ServeHTTP(w, req)
	return w.Result()
}

func (s *APISuiteCore) Get200(c *gc.C, path string) string {
	req := httptest.NewRequest("GET", path, nil)
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 200)
	return readAll(c, res.Body)
}

func (s *APISuiteCore) Get401(c *gc.C, path string) {
	req := httptest.NewRequest("GET", path, nil)
	res := s.Serve(req)
	defer res.Body.Close()
	c.Check(res.StatusCode, gc.Equals, 401)

}

func (s *APISuiteCore) Get400(c *gc.C, path string) {
	req := httptest.NewRequest("GET", path, nil)
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 400)
	res.Body.Close()
}

func (s *APISuiteCore) Get404(c *gc.C, path string) {
	req := httptest.NewRequest("GET", path, nil)
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 404)
	res.Body.Close()
}

func (s *APISuiteCore) Post200(c *gc.C, path, payload string) string {
	req := httptest.NewRequest("POST", path, strings.NewReader(payload))
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 200)
	return readAll(c, res.Body)
}

func (s *APISuiteCore) Post204(c *gc.C, path, payload string) {
	req := httptest.NewRequest("POST", path, strings.NewReader(payload))
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 204)
	res.Body.Close()
}

func (s *APISuiteCore) Post400(c *gc.C, path, payload string) {
	req := httptest.NewRequest("POST", path, strings.NewReader(payload))
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 400)
	res.Body.Close()
}

func (s *APISuiteCore) Post403(c *gc.C, path, payload string) {
	req := httptest.NewRequest("POST", path, strings.NewReader(payload))
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 403)
	res.Body.Close()
}

func (s *APISuiteCore) Post404(c *gc.C, path string) {
	req := httptest.NewRequest("POST", path, nil)
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 404)
	res.Body.Close()
}

func (s *APISuiteCore) Put200(c *gc.C, path, payload string) string {
	req := httptest.NewRequest("PUT", path, strings.NewReader(payload))
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 200)
	return readAll(c, res.Body)
}

func (s *APISuiteCore) Put204(c *gc.C, path, payload string) {
	req := httptest.NewRequest("PUT", path, strings.NewReader(payload))
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 204)
	res.Body.Close()
}

func (s *APISuiteCore) Delete204(c *gc.C, path string) {
	req := httptest.NewRequest("DELETE", path, nil)
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 204)
	res.Body.Close()
}

func (s *APISuiteCore) Delete403(c *gc.C, path string) {
	req := httptest.NewRequest("DELETE", path, nil)
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 403)
	res.Body.Close()
}

func (s *APISuiteCore) Delete404(c *gc.C, path string) {
	req := httptest.NewRequest("DELETE", path, nil)
	res := s.Serve(req)
	c.Check(res.StatusCode, gc.Equals, 404)
	res.Body.Close()
}

func readAll(c *gc.C, rc io.ReadCloser) string {
	defer rc.Close()
	b, err := ioutil.ReadAll(rc)
	c.Assert(err, jc.ErrorIsNil)
	return string(b)
}
