package actions_test

import (
	"github.com/gobuffalo/buffalo/examples/html-crud/models"
	"github.com/markbates/pop"
)

func tx(fn func(tx *pop.Connection)) {
	models.DB.MigrateReset("../migrations")
	fn(models.DB)
}
