package generate

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/buffalo/generators"
	"github.com/markbates/gentronics"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

// ResourceCmd generates a new actions/resource file and a stub test.
var ResourceCmd = &cobra.Command{
	Use:     "resource [name]",
	Aliases: []string{"r"},
	Short:   "Generates a new actions/resource file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you must specify a resource name")
		}
		name := args[0]
		data := gentronics.Data{
			"name":         name,
			"singular":     inflect.Singularize(name),
			"plural":       inflect.Pluralize(name),
			"camel":        inflect.Camelize(name),
			"under":        inflect.Underscore(name),
			"downFirstCap": inflect.CamelizeDownFirst(name),
			"actions":      []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
		}
		return NewResourceGenerator(data).Run(".", data)
	},
}

// NewResourceGenerator generates a new actions/resource file and a stub test.
func NewResourceGenerator(data gentronics.Data) *gentronics.Generator {
	g := gentronics.New()
	g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s.go", data["downFirstCap"])), rAction))
	g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s_test.go", data["under"])), rResourceTest))
	g.Add(&gentronics.Func{
		Should: func(data gentronics.Data) bool { return true },
		Runner: func(root string, data gentronics.Data) error {
			return addInsideAppBlock(fmt.Sprintf("var %sResource buffalo.Resource", data["downFirstCap"]),
				fmt.Sprintf("%sResource = %sResource{&buffalo.BaseResource{}}", data["downFirstCap"], data["camel"]),
				fmt.Sprintf("app.Resource(\"/%s\", %sResource)", data["under"], data["downFirstCap"]),
			)
		},
	})
	g.Add(gentronics.NewCommand(generators.GoFmt()))
	return g
}

var rAction = `package actions

import "github.com/gobuffalo/buffalo"

type {{camel}}Resource struct{
	buffalo.Resource
}

{{#each actions}}
// {{.}} default implementation.
func (v {{camel}}Resource) {{.}}(c buffalo.Context) error {
	return c.Render(200, r.String("{{camel}}#{{.}}"))
}

{{/each}}
`

var rResourceTest = `package actions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)
{{#each actions}}
func Test_{{camel}}Resource_{{camelize .}}(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}
{{/each}}
`
