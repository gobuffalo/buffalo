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
	"fmt"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestGenerateActionArgsComplete(t *testing.T) {
	r := require.New(t)

	cmd := cobra.Command{}

	e := ActionCmd.RunE(&cmd, []string{})
	r.NotNil(e)

	e = ActionCmd.RunE(&cmd, []string{"users"})
	r.NotNil(e)

	os.Chdir(os.TempDir())
	os.Mkdir("actions", 666)

	e = ActionCmd.RunE(&cmd, []string{"users", "show"})
	r.Nil(e)
}

func TestGenerateActionActionsFolderExists(t *testing.T) {
	dir := os.TempDir()

	os.Chdir(dir)
	os.RemoveAll(fmt.Sprintf("%v/actions", dir))

	r := require.New(t)
	cmd := cobra.Command{}

	e := ActionCmd.RunE(&cmd, []string{"users", "show"})
	r.NotNil(e)

	os.Mkdir("actions", 666)

	e = ActionCmd.RunE(&cmd, []string{"users", "show"})
	r.Nil(e)
}
