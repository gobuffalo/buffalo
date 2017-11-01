package action

const (
	actionsHeaderTmpl = `package actions
import "github.com/gobuffalo/buffalo"`

	viewTmpl = `<h1>{{.opts.Name.Camel}}#{{.action}}</h1>`

	actionsTmpl = `
{{ range $action := .actions }}
// {{$.opts.Name.Camel}}{{$action.Camel}} default implementation.
func {{$.opts.Name.Camel}}{{$action.Camel}}(c buffalo.Context) error {
	return c.Render(200, r.HTML("{{$.opts.Name.File}}/{{$action.File}}.html"))
}
{{end}}`

	testHeaderTmpl = `package actions

import (
	"testing"

	"github.com/stretchr/testify/require"
)
	`

	testsTmpl = `
{{ range $action := .tests}}
func (as *ActionSuite) Test_{{$.opts.Name.Camel}}_{{$action.Camel}}() {
	as.Fail("Not Implemented!")
}

{{end}}`
)
