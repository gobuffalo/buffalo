package basicauth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

var (
	// ErrNoCreds is returned when no basic auth credentials are defined
	ErrNoCreds = errors.New("no basic auth credentials defined")

	// ErrAuthFail is returned when the client fails basic authentication
	ErrAuthFail = errors.New("invalid basic auth username or password")
)

// Authorizer is used to authenticate the basic auth username/password.
// Should return true/false and/or an error.
type Authorizer func(buffalo.Context, string, string) (bool, error)

// Middleware enables basic authentication
func Middleware(auth Authorizer) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			token := strings.SplitN(c.Request().Header.Get("Authorization"), " ", 2)
			if len(token) != 2 {
				c.Response().Header().Set("WWW-Authenticate", `Basic realm="Basic Authentication"`)
				return c.Error(http.StatusUnauthorized, errors.New("Unauthorized"))
			}
			b, err := base64.StdEncoding.DecodeString(token[1])
			if err != nil {
				return ErrAuthFail
			}
			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				return ErrAuthFail
			}
			success, err := auth(c, pair[0], pair[1])
			if err != nil {
				return errors.WithStack(err)
			}
			if !success {
				return ErrAuthFail
			}
			return next(c)
		}
	}
}
