package actions

import (
	"os"

	"github.com/markbates/buffalo"
	"github.com/markbates/buffalo/examples/json-resource/models"
	"github.com/markbates/buffalo/middleware"
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

		g := app.Resource("/users", &UsersResource{})
		g.Use(findUserMW)
	}

	return app
}
