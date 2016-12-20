package actions_test

import (
	"github.com/markbates/buffalo/examples/json-crud/models"
	"github.com/markbates/pop"
)

func tx(fn func(tx *pop.Connection)) {
	models.DB.MigrateReset("../migrations")
	fn(models.DB)
}
