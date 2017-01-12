package generate

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
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

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
		app = buffalo.Automatic(buffalo.Options{
			Env: ENV,
		})

		app.Use(middleware.SetContentType("application/json"))
		app.Use(middleware.PopTransaction(models.DB))

		app.Use(findUserMW)
		app.GET("/users", UsersList)
	}

	return app
}`

	ioutil.WriteFile(filepath.Join(packagePath, "actions", "app.go"), []byte(shortAppFileExample), 0755)

	addRoute("GET", "/new/route", "UserCoolHandler")

	contentAfter, _ := ioutil.ReadFile(filepath.Join(packagePath, "actions", "app.go"))
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
		app = buffalo.Automatic(buffalo.Options{
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
