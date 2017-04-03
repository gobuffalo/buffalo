package generate

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestGenerateResourceCode(t *testing.T) {
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

	SkipResourceMigration = false
	SkipResourceModel = false

	// Testing generator with singular definition
	e = ResourceCmd.RunE(&cmd, []string{"user"})
	r.Nil(e)

	fileData, _ := ioutil.ReadFile("actions/app.go")
	r.Contains(string(fileData), "var usersResource buffalo.Resource")
	r.Contains(string(fileData), "usersResource = UsersResource{&buffalo.BaseResource{}}")
	r.Contains(string(fileData), "app.Resource(\"/users\", usersResource)")

	fileData, _ = ioutil.ReadFile("actions/users.go")
	r.Contains(string(fileData), "type UsersResource struct {")
	r.Contains(string(fileData), "func (v UsersResource) List(c buffalo.Context) error {")
	r.Contains(string(fileData), "func (v UsersResource) Destroy(c buffalo.Context) error {")

	fileData, _ = ioutil.ReadFile("actions/users_test.go")
	r.Contains(string(fileData), "func (as *ActionSuite) Test_UsersResource_List")
	r.Contains(string(fileData), "func (as *ActionSuite) Test_UsersResource_Show")
	r.Contains(string(fileData), "func (as *ActionSuite) Test_UsersResource_Create")

	// Testing generator with plural definition
	e = ResourceCmd.RunE(&cmd, []string{"comments"})
	r.Nil(e)

	fileData, _ = ioutil.ReadFile("actions/app.go")
	r.Contains(string(fileData), "var commentsResource buffalo.Resource")
	r.Contains(string(fileData), "commentsResource = CommentsResource{&buffalo.BaseResource{}}")
	r.Contains(string(fileData), "app.Resource(\"/comments\", commentsResource)")

	fileData, _ = ioutil.ReadFile("actions/comments.go")
	r.Contains(string(fileData), "type CommentsResource struct {")
	r.Contains(string(fileData), "func (v CommentsResource) List(c buffalo.Context) error {")
	r.Contains(string(fileData), "func (v CommentsResource) Destroy(c buffalo.Context) error {")

	fileData, _ = ioutil.ReadFile("actions/comments_test.go")
	r.Contains(string(fileData), "func (as *ActionSuite) Test_CommentsResource_List")
	r.Contains(string(fileData), "func (as *ActionSuite) Test_CommentsResource_Show")
	r.Contains(string(fileData), "func (as *ActionSuite) Test_CommentsResource_Create")

}
