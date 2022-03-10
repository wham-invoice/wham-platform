package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

// HTTPError can be returned by an Endpoint.Do func to signal an HTTP status.
type HTTPError struct {
	Status int
}

// Error is part of the error interface.
func (e HTTPError) Error() string {
	return fmt.Sprintf("we should return an HTTP %d", e.Status)
}

// Abort aborts the context with the supplied error, setting the response
// status if err is an HTTPError and otherwise setting 500 an logging. It will
// panic if passed a nil error.
func Abort(c *gin.Context, err error) {
	switch terr := errors.Cause(err).(type) {
	case nil:
		panic("Abort with nil error")
	case HTTPError:
		c.AbortWithStatus(terr.Status)

	default:

		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

var (
	BadRequest   = HTTPError{http.StatusBadRequest}
	Unauthorized = HTTPError{http.StatusUnauthorized}
	NotFound     = HTTPError{http.StatusNotFound}
	Forbidden    = HTTPError{http.StatusForbidden}
)
