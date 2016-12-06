package actions_test

import (
	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/html-crud/models"
	"github.com/markbates/buffalo/middleware"
	"github.com/markbates/pop"
)

func tx(fn func(tx *pop.Connection)) {
	tmw := middleware.PopTransaction
	defer func() {
		middleware.PopTransaction = tmw
	}()

	models.DB.MigrateReset("../migrations")
	models.DB.Rollback(func(tx *pop.Connection) {
		middleware.PopTransaction = func(db *pop.Connection) buffalo.MiddlewareFunc {
			return func(h buffalo.Handler) buffalo.Handler {
				return func(c buffalo.Context) error {
					c.Set("tx", tx)
					return h(c)
				}
			}
		}
		fn(tx)
	})
}
