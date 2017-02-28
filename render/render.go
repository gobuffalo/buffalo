package render

import (
	"bytes"
	"html/template"

	"github.com/gobuffalo/buffalo/render/resolvers"
	"github.com/gobuffalo/velvet"
)

// Engine used to power all defined renderers.
// This allows you to configure the system to your
// preferred settings, instead of just getting
// the defaults.
type Engine struct {
	Options
}

// New render.Engine ready to go with your Options
// and some defaults we think you might like.
func New(opts Options) *Engine {
	if opts.Helpers == nil {
		opts.Helpers = map[string]interface{}{}
	}
	if opts.FileResolverFunc == nil {
		opts.FileResolverFunc = func() resolvers.FileResolver {
			return &resolvers.SimpleResolver{}
		}
	}
	if opts.TemplateEngine == nil {
		opts.TemplateEngine = velvet.BuffaloRenderer
	}

	e := &Engine{
		Options: opts,
	}
	return e
}

type TemplateEngine func(string, map[string]interface{}, map[string]interface{}) (string, error)

func GoTemplateEngine(input string, data map[string]interface{}, helpers map[string]interface{}) (string, error) {
	t, err := template.New(input).Parse(input)
	if err != nil {
		return "", err
	}
	if helpers != nil {
		t = t.Funcs(helpers)
	}
	bb := &bytes.Buffer{}
	err = t.Execute(bb, data)
	return bb.String(), err
}
