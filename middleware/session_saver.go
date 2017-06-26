package middleware

import (
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// SessionSaver will automatically save a session if the
// request was successful.
func SessionSaver(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		err := next(c)
		if err != nil {
			return errors.WithStack(err)
		}
		return c.Session().Save()
	}
}
