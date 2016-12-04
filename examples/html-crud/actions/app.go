package actions

import (
	"net/http"

	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/html-crud/models"
)

func App() http.Handler {
	a := buffalo.Automatic(buffalo.Options{})
	a.Env = "development"

	a.ServeFiles("/assets", assetsPath())
	a.Use(models.TransactionMW)
	a.GET("/", func(c buffalo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/users")
	})

	a.GET("/user/new", UsersNew)

	g := a.Group("/users")
	g.Use(findUserMW)
	g.GET("/", UsersList)
	g.GET("/:user_id", UsersShow)
	g.GET("/:user_id/edit", UsersEdit)
	g.POST("/", UsersCreate)
	g.PUT("/:user_id", UsersUpdate)
	g.DELETE("/:user_id", UsersDelete)

	return a
}
