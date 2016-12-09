package actions

import (
	"net/http"

	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/json-crud/models"
	"github.com/markbates/buffalo/middleware"
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
	a.Use(middleware.PopTransaction(models.DB))

	a.Use(findUserMW)
	a.GET("/users", UsersList)
	a.GET("/users/{user_id}", UsersShow)
	a.POST("/users", UsersCreate)
	a.PUT("/users/{user_id}", UsersUpdate)
	a.DELETE("/users/{user_id}", UsersDelete)
	return a
}
