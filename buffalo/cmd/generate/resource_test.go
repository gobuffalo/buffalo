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

	e = ResourceCmd.RunE(&cmd, []string{"users"})
	r.Nil(e)

	fileData, _ := ioutil.ReadFile("actions/app.go")
	r.Contains(string(fileData), "var usersResource buffalo.Resource")
	r.Contains(string(fileData), "usersResource = &UsersResource{&buffalo.BaseResource{}}")
	r.Contains(string(fileData), "app.Resource(\"/users\", usersResource)")

}
