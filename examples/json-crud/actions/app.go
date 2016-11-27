package actions

import (
	"net/http"
	"time"

	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/json-crud/models"
	"github.com/markbates/pop"
)

func App() http.Handler {
	a := buffalo.Automatic(buffalo.Options{})
	a.Use(func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// since this is a JSON-only app, let's make sure the
			// content-type is always json
			c.Request().Header.Set("Content-Type", "application/json")
			return h(c)
		}
	})
	a.Use(func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// wrap all requests in a transaction and set the length
			// of time doing things in the db to the log.
			return models.DB.Transaction(func(tx *pop.Connection) error {
				start := tx.Elapsed
				defer func() {
					finished := tx.Elapsed
					elapsed := time.Duration(finished - start)
					c.LogField("db", elapsed)
				}()
				c.Set("tx", tx)
				return h(c)
			})
		}
	})

	a.GET("/users", UsersList)
	a.GET("/users/:user_id", UsersShow)
	a.POST("/users", UsersCreate)
	a.PUT("/users/:user_id", UsersUpdate)
	a.DELETE("/users/:user_id", UsersDelete)
	return a
}
