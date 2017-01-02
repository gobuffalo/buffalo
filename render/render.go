package render

import (
	"sync"

	"github.com/gobuffalo/buffalo/render/resolvers"
	"github.com/gobuffalo/velvet"
)

// Engine used to power all defined renderers.
// This allows you to configure the system to your
// prefered settings, instead of just getting
// the defaults.
type Engine struct {
	Options
	templateCache map[string]*velvet.Template
	moot          *sync.Mutex
}

// New render.Engine ready to go with your Options
// and some defaults we think you might like. Engines
// have the following helpers added to them:
// https://github.com/gobuffalo/buffalo/blob/master/render/helpers/helpers.go#L1
// https://github.com/markbates/inflect/blob/master/helpers.go#L3
func New(opts Options) *Engine {
	if opts.FileResolver == nil {
		opts.FileResolver = &resolvers.SimpleResolver{}
	}

	e := &Engine{
		Options:       opts,
		templateCache: map[string]*velvet.Template{},
		moot:          &sync.Mutex{},
	}
	return e
}

// RegisterHelper adds a helper to a template with the given name.
// See github.com/gobuffalo/velvet for more details on helpers.
/*
	e.RegisterHelper("upcase", strings.ToUpper)
*/
func (e *Engine) RegisterHelper(name string, helper interface{}) error {
	return velvet.Helpers.Add(name, helper)
}

// RegisterHelpers adds helpers to a template with the given name.
// See github.com/gobuffalo/velvet for more details on helpers.
/*
	h := map[string]interface{}{
		"upcase": strings.ToUpper,
		"downcase": strings.ToLower,
	}
	e.RegisterHelpers(h)
*/
func (e *Engine) RegisterHelpers(helpers map[string]interface{}) error {
	for k, v := range helpers {
		err := e.RegisterHelper(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
