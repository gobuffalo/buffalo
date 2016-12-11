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
			return errors.New("You must specifiy a resource name!")
		}
		name := args[0]
		data := gentronics.Data{
			"name":       name,
			"singular":   inflect.Singularize(name),
			"plural":     inflect.Pluralize(name),
			"camel":      inflect.Camelize(name),
			"underscore": inflect.Underscore(name),
		}
		return NewResourceGenerator(data).Run(".", data)
	},
}

// NewResourceGenerator generates a new actions/resource file and a stub test.
func NewResourceGenerator(data gentronics.Data) *gentronics.Generator {
	g := gentronics.New()
	g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s.go", data["underscore"])), rAction))
	g.Add(gentronics.NewFile(filepath.Join("actions", fmt.Sprintf("%s_test.go", data["underscore"])), rActionTest))
	g.Add(Fmt)
	return g
}

var rAction = `package actions

import "github.com/markbates/buffalo"

type {{.camel}}Resource struct{}

// List default implementation. Returns a 404
func (v *{{.camel}}Resource) List(c buffalo.Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

// Show default implementation. Returns a 404
func (v *{{.camel}}Resource) Show(c buffalo.Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

// New default implementation. Returns a 404
func (v *{{.camel}}Resource) New(c buffalo.Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

// Create default implementation. Returns a 404
func (v *{{.camel}}Resource) Create(c buffalo.Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

// Edit default implementation. Returns a 404
func (v *{{.camel}}Resource) Edit(c buffalo.Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

// Update default implementation. Returns a 404
func (v *{{.camel}}Resource) Update(c buffalo.Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}

// Destroy default implementation. Returns a 404
func (v *{{.camel}}Resource) Destroy(c buffalo.Context) error {
	return c.Error(404, errors.New("resource not implemented"))
}`

var rActionTest = `package actions_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_{{.camel}}Resource_List(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{.camel}}Resource_Show(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{.camel}}Resource_New(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{.camel}}Resource_Create(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{.camel}}Resource_Edit(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{.camel}}Resource_Update(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}

func Test_{{.camel}}Resource_Destroy(t *testing.T) {
	r := require.New(t)
	r.Fail("Not Implemented!")
}
`
