package generators

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendRoute(t *testing.T) {
	r := require.New(t)

	tmpDir := os.TempDir()
	packagePath := filepath.Join(tmpDir, "src", "sample")
	err := os.MkdirAll(packagePath, 0777)
	r.NoError(err)

	err = os.Chdir(packagePath)
	r.NoError(err)

	const shortAppFileExample = `package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/examples/json-crud/models"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/markbates/going/defaults"
)

var ENV = defaults.String(os.Getenv("GO_ENV"), "development")
var app *buffalo.App
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env: ENV,
		})

		app.Use(middleware.SetContentType("application/json"))
		app.Use(middleware.PopTransaction(models.DB))

		app.Use(findUserMW)
		app.GET("/users", UsersList)
	}

	return app
}`

	err = os.MkdirAll("actions", 0777)
	r.NoError(err)
	err = ioutil.WriteFile(filepath.Join(packagePath, "actions", "app.go"), []byte(shortAppFileExample), 0755)
	r.NoError(err)

	err = AddRoute("GET", "/new/route", "UserCoolHandler")
	r.NoError(err)

	contentAfter, err := ioutil.ReadFile(filepath.Join(packagePath, "actions", "app.go"))
	r.NoError(err)
	r.Equal(`package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/examples/json-crud/models"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/markbates/going/defaults"
)

var ENV = defaults.String(os.Getenv("GO_ENV"), "development")
var app *buffalo.App
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env: ENV,
		})

		app.Use(middleware.SetContentType("application/json"))
		app.Use(middleware.PopTransaction(models.DB))

		app.Use(findUserMW)
		app.GET("/users", UsersList)
		app.GET("/new/route", UserCoolHandler)
	}

	return app
}`, string(contentAfter))

}
