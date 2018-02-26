package middleware

import (
	"time"

	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
)

// PopTransaction is a piece of Buffalo middleware that wraps each
// request in a transaction that will automatically get committed or
// rolledback. It will also add a field to the log, "db", that
// shows the total duration spent during the request making database
// calls.
var PopTransaction = func(db *pop.Connection) buffalo.MiddlewareFunc {
	return func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// wrap all requests in a transaction and set the length
			// of time doing things in the db to the log.
			err := db.Transaction(func(tx *pop.Connection) error {
				start := tx.Elapsed
				defer func() {
					finished := tx.Elapsed
					elapsed := time.Duration(finished - start)
					c.LogField("db", elapsed)
				}()
				c.Set("tx", tx)
				if err := h(c); err != nil {
					return err
				}
				if res, ok := c.Response().(*buffalo.Response); ok {
					if res.Status < 200 || res.Status >= 400 {
						return errNonSuccess
					}
				}
				return nil
			})
			if err != nil && errors.Cause(err) != errNonSuccess {
				return err
			}
			return nil
		}
	}
}

var errNonSuccess = errors.New("non success status code")
