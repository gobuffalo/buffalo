package grifts

import (
	"github.com/gobuffalo/buffalo/examples/html-crud/models"
	"github.com/markbates/grift/grift"
	"github.com/markbates/pop"
)

var _ = grift.Add("db:seed", func(c *grift.Context) error {
	return models.DB.Transaction(func(tx *pop.Connection) error {
		users := models.Users{
			{FirstName: "Mark", LastName: "Bates", Email: "mark@example.com"},
			{FirstName: "Jane", LastName: "Doe", Email: "jane@example.com"},
		}
		for _, u := range users {
			err := tx.Create(&u)
			if err != nil {
				return err
			}
		}
		return nil
	})
})
