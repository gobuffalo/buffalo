package generate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestGenerateActionArgsComplete(t *testing.T) {
	dir := os.TempDir()
	packagePath := filepath.Join(dir, "src", "sample")
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	r := require.New(t)

	cmd := cobra.Command{}

	e := ActionCmd.RunE(&cmd, []string{})
	r.NotNil(e)

	e = ActionCmd.RunE(&cmd, []string{"users"})
	r.NotNil(e)

	os.Mkdir("actions", 0755)
	ioutil.WriteFile("actions/app.go", appGo, 0755)

	e = ActionCmd.RunE(&cmd, []string{"users", "show"})
	r.Nil(e)
}

func TestGenerateActionActionsFolderExists(t *testing.T) {
	dir := os.TempDir()
	packagePath := filepath.Join(dir, "src", "sample")
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	os.RemoveAll("actions")
	os.RemoveAll("templates")

	r := require.New(t)
	cmd := cobra.Command{}

	e := ActionCmd.RunE(&cmd, []string{"users", "show", "edit"})
	r.NotNil(e)

	os.Mkdir("actions", 0755)
	ioutil.WriteFile("actions/app.go", appGo, 0755)

	e = ActionCmd.RunE(&cmd, []string{"users", "show", "edit"})
	r.Nil(e)

	data, _ := ioutil.ReadFile("actions/users.go")
	r.Contains(string(data), "package actions")
	r.Contains(string(data), `import "github.com/gobuffalo/buffalo"`)
	r.Contains(string(data), "func UsersShow(c buffalo.Context) error {")
	r.Contains(string(data), "func UsersEdit(c buffalo.Context) error {")
	r.Contains(string(data), `r.HTML("users/edit.html")`)
	r.Contains(string(data), `c.Render(200, r.HTML("users/show.html"))`)

	data, _ = ioutil.ReadFile("templates/users/show.html")
	r.Contains(string(data), "<h1>Users#Show</h1>")
}

func TestGenerateActionActionsFileExists(t *testing.T) {
	dir := os.TempDir()
	packagePath := filepath.Join(dir, "src", "sample")
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	os.Mkdir("actions", 0755)
	ioutil.WriteFile("actions/app.go", appGo, 0755)
	r := require.New(t)
	cmd := cobra.Command{}
	usersContent := `package actions
import "log"

func UsersShow(c buffalo.Context) error {
    log.Println("Something Here!")
    return c.Render(200, r.String("OK"))
}
`
	ioutil.WriteFile("actions/users.go", []byte(usersContent), 0755)

	e := ActionCmd.RunE(&cmd, []string{"users", "show"})
	r.Nil(e)

	data, _ := ioutil.ReadFile("actions/users.go")
	r.Contains(string(data), "log.Println(")
	r.Contains(string(data), "func UsersShow")

}

func TestGenerateNewActionWithExistingActions(t *testing.T) {
	dir := os.TempDir()
	packagePath := filepath.Join(dir, "src", "sample")
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	os.RemoveAll("actions")
	os.RemoveAll("templates")

	os.Mkdir("actions", 0755)
	ioutil.WriteFile("actions/app.go", appGo, 0755)
	r := require.New(t)
	cmd := cobra.Command{}
	e := ActionCmd.RunE(&cmd, []string{"users", "show", "edit"})
	r.Nil(e)

	data, _ := ioutil.ReadFile("actions/users.go")
	r.Contains(string(data), "package actions")
	r.Contains(string(data), `import "github.com/gobuffalo/buffalo"`)
	r.Contains(string(data), "func UsersShow(c buffalo.Context) error {")
	r.Contains(string(data), "func UsersEdit(c buffalo.Context) error {")
	r.Contains(string(data), `r.HTML("users/edit.html")`)
	r.Contains(string(data), `c.Render(200, r.HTML("users/show.html"))`)

	e = ActionCmd.RunE(&cmd, []string{"users", "list"})
	r.Nil(e)

	data, _ = ioutil.ReadFile("actions/users.go")
	r.Contains(string(data), "package actions")
	r.Contains(string(data), `import "github.com/gobuffalo/buffalo"`)
	r.Contains(string(data), "func UsersShow(c buffalo.Context) error {")
	r.Contains(string(data), "func UsersEdit(c buffalo.Context) error {")
	r.Contains(string(data), "func UsersList(c buffalo.Context) error {")
	r.Contains(string(data), `r.HTML("users/list.html")`)
	r.Contains(string(data), `c.Render(200, r.HTML("users/list.html"))`)

	data, _ = ioutil.ReadFile("templates/users/list.html")
	r.Contains(string(data), "<h1>Users#List</h1>")

	data, _ = ioutil.ReadFile("actions/users_test.go")
	r.Contains(string(data), "package actions_test")
	r.Contains(string(data), "func (as *ActionSuite) Test_Users_Show() {")
	r.Contains(string(data), "func (as *ActionSuite) Test_Users_Edit() {")
	r.Contains(string(data), "func (as *ActionSuite) Test_Users_List() {")
}

var appGo = []byte(`
package actions
var app *buffalo.App
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env: "test",
		})
		app.GET("/", func (c buffalo.Context) error {
			return c.Render(200, r.String("hello"))
		})
	}

	return app
}
`)
