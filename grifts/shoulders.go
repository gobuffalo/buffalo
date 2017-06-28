// +build !appengine

package grifts

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/markbates/deplist"
	"github.com/markbates/grift/grift"
)

var _ = grift.Desc("shoulders", "Prints a listing all of the 3rd party packages used by buffalo.")
var _ = grift.Add("shoulders:list", func(c *grift.Context) error {
	giants, _ := deplist.List("examples")
	for _, k := range []string{
		"github.com/markbates/refresh",
		"github.com/markbates/grift",
		"github.com/markbates/pop",
		"github.com/spf13/cobra",
		"github.com/motemen/gore",
		"golang.org/x/tools/cmd/goimports",
	} {
		giants[k] = k
	}

	deps := make([]string, 0, len(giants))
	for k := range giants {
		if !strings.Contains(k, "github.com/gobuffalo/buffalo") {
			deps = append(deps, k)
		}
	}
	sort.Strings(deps)
	fmt.Println(strings.Join(deps, "\n"))
	c.Set("giants", deps)
	return nil
})

var _ = grift.Desc("shoulders", "Generates a file listing all of the 3rd party packages used by buffalo.")
var _ = grift.Add("shoulders", func(c *grift.Context) error {
	err := grift.Run("shoulders:list", c)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(envy.GoPath(), "src", "github.com", "gobuffalo", "buffalo", "SHOULDERS.md"))
	if err != nil {
		return err
	}
	t, err := template.New("").Parse(shouldersTemplate)
	if err != nil {
		return err
	}
	err = t.Execute(f, c.Value("giants"))
	if err != nil {
		return err
	}

	return commitAndPushShoulders()
})

func commitAndPushShoulders() error {
	cmd := exec.Command("git", "commit", "SHOULDERS.md", "-m", "Updated SHOULDERS.md")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "push", "origin")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

var shouldersTemplate = `
# Buffalo Stands on the Shoulders of Giants

Buffalo does not try to reinvent the wheel! Instead, it uses the already great wheels developed by the Go community and puts them altogether in the best way possible. Without these giants this project would not be possible. Please make sure to check them out and thank them for all of their hard work.

Thank you to the following **GIANTS**:

{{ range $v := .}}
* [{{$v}}](https://{{$v}})
{{ end }}
`
