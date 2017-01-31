// Copyright © 2016 Mark Bates <mark@markbates.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package generate

import (
	"errors"
	"fmt"
	"path/filepath"

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
			"name":     name,
			"singular": inflect.Singularize(name),
			"plural":   inflect.Pluralize(name),
			"camel":    inflect.Camelize(name),
			"under":    inflect.Underscore(name),
			"actions":  []string{"List", "Show", "New", "Create", "Edit", "Update", "Destroy"},
		}
		return NewResourceGenerator(data).Run(".", data)
	},
}

// NewResourceGenerator generates a new actions/resource file and a stub test.
func NewResourceGenerator(data gentronics.Data) *gentronics.Generator {
	g := gentronics.New()
	g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s.go", data["under"])), rAction))
	g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s_test.go", data["under"])), rResourceTest))
	g.Add(&gentronics.Func{
		Should: func(data gentronics.Data) bool { return true },
		Runner: func(root string, data gentronics.Data) error {
			return addInsideAppBlock(fmt.Sprintf("var %sresource buffalo.Resource", data["camel"]),
				fmt.Sprintf("%sresource = &%sResource{&buffalo.BaseResource{}}", data["camel"], data["camel"]),
				fmt.Sprintf("app.Resource(\"/%s\", %sresource)", data["under"], data["camel"]),
			)
		},
	})
	g.Add(Fmt)
	return g
}

var rAction = `package actions

import "github.com/gobuffalo/buffalo"

type {{camel}}Resource struct{
	buffalo.Resource
}

{{#each actions}}
// {{.}} default implementation.
func (v *{{camel}}Resource) {{.}}(c buffalo.Context) error {
	return c.Render(200, r.String("{{camel}}#{{.}}"))
}

{{/each}}
`

var rResourceTest = `package actions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_{{camel}}Resource_List(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{camel}}Resource_Show(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{camel}}Resource_New(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{camel}}Resource_Create(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{camel}}Resource_Edit(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{camel}}Resource_Update(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{camel}}Resource_Destroy(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}
`
