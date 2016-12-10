package actions

import (
	"net/http"

	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/html-resource/models"
	"github.com/markbates/buffalo/middleware"
)

func App() http.Handler {
	a := buffalo.Automatic(buffalo.Options{})
	a.Env = "development"

	a.ServeFiles("/assets", assetsPath())
	a.Use(middleware.PopTransaction(models.DB))
	a.GET("/", func(c buffalo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/users")
	})

	a.Use(findUserMW("user_id"))
	a.Resource("/users", &UsersResource{})
	a.Use(findUserMW("person_id"))
	a.Resource("/people", &UsersResource{})

	return a
}
