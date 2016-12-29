package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/examples/json-crud/models"
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

		app.Use(middleware.SetContentType("application/json"))
		app.Use(middleware.PopTransaction(models.DB))

		app.Use(findUserMW)
		app.GET("/users", UsersList)
		app.GET("/users/{user_id}", UsersShow)
		app.POST("/users", UsersCreate)
		app.PUT("/users/{user_id}", UsersUpdate)
		app.DELETE("/users/{user_id}", UsersDelete)
	}

	return app
}
