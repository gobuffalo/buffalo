package middleware

import (
	"fmt"

	"github.com/markbates/buffalo"
	newrelic "github.com/newrelic/go-agent"
	"github.com/pkg/errors"
)

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
