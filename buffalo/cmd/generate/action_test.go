// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	e = ActionCmd.RunE(&cmd, []string{"users", "show", "edit"})
	r.Nil(e)

	data, _ := ioutil.ReadFile("actions/users.go")
	r.Contains(string(data), "func UsersShow(c buffalo.Context) error {")
	r.Contains(string(data), "func UsersEdit(c buffalo.Context) error {")

	data, _ = ioutil.ReadFile("templates/users/show.html")
	r.Contains(string(data), "<h1>Users#Show</h1>")
}

func TestGenerateActionActionsFileExists(t *testing.T) {
	dir := os.TempDir()
	packagePath := filepath.Join(dir, "src", "sample")
	os.MkdirAll(packagePath, 0755)
	os.Chdir(packagePath)

	os.Mkdir("actions", 0755)
	r := require.New(t)
	cmd := cobra.Command{}
	usersContent := `package actions
import log

func UsersShow(c buffalo.Context)}{
    log.Println("Something Here!")
}
`
	ioutil.WriteFile("actions/users.go", []byte(usersContent), 0755)

	e := ActionCmd.RunE(&cmd, []string{"users", "show", "edit", "other"})
	r.Nil(e)

	data, _ := ioutil.ReadFile("actions/users.go")
	r.Contains(string(data), "log.Println(")
	r.Contains(string(data), "func UsersEdit")
	r.Contains(string(data), "func UsersOther")

}
