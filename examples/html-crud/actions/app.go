package actions

import (
	"net/http"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/examples/html-crud/models"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/markbates/going/defaults"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = defaults.String(os.Getenv("GO_ENV"), "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env: ENV,
		})

		app.ServeFiles("/assets", assetsPath())
		app.Use(middleware.PopTransaction(models.DB))
		app.GET("/", func(c buffalo.Context) error {
			return c.Redirect(http.StatusPermanentRedirect, "/users")
		})

		g := app.Group("/users")
		g.Use(findUserMW)
		g.GET("/", UsersList)
		g.GET("/new", UsersNew)
		g.GET("/{user_id}", UsersShow)
		g.GET("/{user_id}/edit", UsersEdit)
		g.POST("/", UsersCreate)
		g.PUT("/{user_id}", UsersUpdate)
		g.DELETE("/{user_id}", UsersDelete)
	}

	return app
}
