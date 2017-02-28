package render

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/gobuffalo/buffalo/render/resolvers"
)

// Engine used to power all defined renderers.
// This allows you to configure the system to your
// preferred settings, instead of just getting
// the defaults.
type Engine struct {
	Options
	moot *sync.Mutex
}

// New render.Engine ready to go with your Options
// and some defaults we think you might like. Engines
// have the following helpers added to them:
// https://github.com/gobuffalo/buffalo/blob/master/render/helpers/helpers.go#L1
// https://github.com/markbates/inflect/blob/master/helpers.go#L3
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
		opts.TemplateEngine = GoTemplateEngine
	}

	e := &Engine{
		Options: opts,
		moot:    &sync.Mutex{},
	}
	return e
}

type TemplateOptions struct {
	Data    map[string]interface{}
	Helpers map[string]interface{}
}

type TemplateEngine func(string, TemplateOptions) (string, error)

func GoTemplateEngine(input string, opts TemplateOptions) (string, error) {
	t, err := template.New(input).Parse(input)
	if err != nil {
		return "", err
	}
	t = t.Funcs(opts.Helpers)
	bb := &bytes.Buffer{}
	err = t.Execute(bb, opts.Data)
	return bb.String(), err
}
