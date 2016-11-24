package actions_test

import (
	"log"

	"github.com/markbates/buffalo/examples/json-crud/models"
	"github.com/markbates/pop"
)

func tx(fn func()) {
	odb := models.DB
	defer func() {
		models.DB = odb
	}()
	var err error
	models.DB, err = pop.Connect("test")
	if err != nil {
		log.Fatal(err)
	}
	models.DB.Rollback(func(tx *pop.Connection) {
		models.DB = tx
	})
}
