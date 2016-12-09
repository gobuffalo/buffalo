package actions

import (
	"net/http"

	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/html-crud/models"
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

	g := a.Group("/users")
	g.Use(findUserMW)
	g.GET("/", UsersList)
	g.GET("/new", UsersNew)
	g.GET("/{user_id}", UsersShow)
	g.GET("/{user_id}/edit", UsersEdit)
	g.POST("/", UsersCreate)
	g.PUT("/{user_id}", UsersUpdate)
	g.DELETE("/{user_id}", UsersDelete)

	return a
}
