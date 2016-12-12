package middleware

import (
	"fmt"

	"github.com/markbates/buffalo"
	newrelic "github.com/newrelic/go-agent"
	"github.com/pkg/errors"
)

// NewRelic returns a piece of buffalo.Middleware that can
// be used to report requests to NewRelic. You must pass in your
// NewRelic key and a name for your application. If the key
// passed in is blank, i.e. loading from an ENV, then the middleware
// is skipped and the chain continues on like normal. Useful
// for development.
func NewRelic(key, name string) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		if key == "" {
			return next
		}
		return func(c buffalo.Context) error {
			fmt.Printf("Setting up New Relic %s\n", key)
			config := newrelic.NewConfig(name, key)
			app, err := newrelic.NewApplication(config)
			if err != nil {
				return errors.WithStack(err)
			}
			tx := app.StartTransaction(c.Request().URL.String(), c.Response(), c.Request())
			defer tx.End()
			return next(c)
		}
	}
}
