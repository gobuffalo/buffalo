package action

const (
	rActionFileT = `package actions
import "github.com/gobuffalo/buffalo"`

	rViewT = `<h1>{{.namespace}}#{{.action}}</h1>`

	rActionFuncT = `
// {{.namespace}}{{.action}} default implementation.
func {{.namespace}}{{.action}}(c buffalo.Context) error {
	return c.Render(200, r.HTML("{{.namespace_under}}/{{.action_under}}.html"))
}
`
)
