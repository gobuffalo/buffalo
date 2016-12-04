package actions_test

import (
	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/html-crud/models"
	"github.com/markbates/pop"
)

func tx(fn func(tx *pop.Connection)) {
	tmw := models.TransactionMW
	defer func() {
		models.TransactionMW = tmw
	}()

	models.DB.MigrateReset("../migrations")
	models.DB.Rollback(func(tx *pop.Connection) {
		models.TransactionMW = func(h buffalo.Handler) buffalo.Handler {
			return func(c buffalo.Context) error {
				c.Set("tx", tx)
				return h(c)
			}
		}
		fn(tx)
	})
}
