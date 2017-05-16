package generate

// func TestGenerateActionArgsComplete(t *testing.T) {
// 	dir := os.TempDir()
// 	packagePath := filepath.Join(dir, "src", "sample")
// 	os.MkdirAll(packagePath, 0755)
// 	os.Chdir(packagePath)

// 	r := require.New(t)

// 	cmd := cobra.Command{}

// 	e := ActionCmd.RunE(&cmd, []string{})
// 	r.NotNil(e)

// 	e = ActionCmd.RunE(&cmd, []string{"nodes"})
// 	r.NotNil(e)

// 	os.Mkdir("actions", 0755)
// 	ioutil.WriteFile("actions/app.go", appGo, 0755)

// 	e = ActionCmd.RunE(&cmd, []string{"nodes", "show"})
// 	r.Nil(e)
// }

/*
func TestGenerateActionActionsFolderExists(t *testing.T) {
	dir := os.TempDir()
	packagePath := filepath.Join(dir, "src", "sample")
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	os.RemoveAll("actions")
	os.RemoveAll("templates")

	r := require.New(t)
	cmd := cobra.Command{}

	e := ActionCmd.RunE(&cmd, []string{"comments", "show", "edit"})
	r.NotNil(e)

	os.Mkdir("actions", 0755)
	ioutil.WriteFile("actions/app.go", appGo, 0755)

	e = ActionCmd.RunE(&cmd, []string{"comments", "show", "edit"})
	r.Nil(e)

	data, _ := ioutil.ReadFile("actions/comments.go")
	r.Contains(string(data), "package actions")
	r.Contains(string(data), `import "github.com/gobuffalo/buffalo"`)
	r.Contains(string(data), "func CommentsShow(c buffalo.Context) error {")
	r.Contains(string(data), "func CommentsEdit(c buffalo.Context) error {")
	r.Contains(string(data), `r.HTML("comments/edit.html")`)
	r.Contains(string(data), `c.Render(200, r.HTML("comments/show.html"))`)

	data, _ = ioutil.ReadFile("templates/comments/show.html")
	r.Contains(string(data), "")
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
	e := ActionCmd.RunE(&cmd, []string{"posts", "show", "edit"})
	r.Nil(e)

	data, _ := ioutil.ReadFile("actions/posts.go")
	r.Contains(string(data), "package actions")
	r.Contains(string(data), "github.com/gobuffalo/buffalo")
	r.Contains(string(data), "func PostsShow(c buffalo.Context) error {")
	r.Contains(string(data), "func PostsEdit(c buffalo.Context) error {")
	r.Contains(string(data), `r.HTML("posts/edit.html")`)
	r.Contains(string(data), `c.Render(200, r.HTML("posts/show.html"))`)

	e = ActionCmd.RunE(&cmd, []string{"posts", "list"})
	r.Nil(e)

	data, _ = ioutil.ReadFile("actions/posts.go")
	r.Contains(string(data), "package actions")
	r.Contains(string(data), "github.com/gobuffalo/buffalo")
	r.Contains(string(data), "func PostsShow(c buffalo.Context) error {")
	r.Contains(string(data), "func PostsEdit(c buffalo.Context) error {")
	r.Contains(string(data), "func PostsList(c buffalo.Context) error {")
	r.Contains(string(data), `c.Render(200, r.HTML("posts/list.html"))`)

	data, _ = ioutil.ReadFile("templates/posts/list.html")
	r.Contains(string(data), "<h1>Posts#List</h1>")

	data, _ = ioutil.ReadFile("actions/posts_test.go")
	r.Contains(string(data), "package actions_test")
	r.Contains(string(data), "func (as *ActionSuite) Test_Posts_Show() {")
	r.Contains(string(data), "func (as *ActionSuite) Test_Posts_Edit() {")
	r.Contains(string(data), "func (as *ActionSuite) Test_Posts_List() {")

	e = ActionCmd.RunE(&cmd, []string{"posts", "list"})
	r.Nil(e)

	data, _ = ioutil.ReadFile("actions/posts_test.go")
	r.Equal(strings.Count(string(data), "func (as *ActionSuite) Test_Posts_List() {"), 1)

	data, _ = ioutil.ReadFile("actions/app.go")
	r.Equal(strings.Count(string(data), "app.GET(\"/posts/list\", PostsList)"), 1)
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
*/
