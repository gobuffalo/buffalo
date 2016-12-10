package grifts

import (
	"github.com/markbates/buffalo/examples/html-resource/models"
	. "github.com/markbates/grift/grift"
	"github.com/markbates/pop"
)

var _ = Add("db:seed", func(c *Context) error {
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
