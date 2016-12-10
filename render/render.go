package render

import (
	"github.com/aymerick/raymond"
	"github.com/markbates/buffalo/render/helpers"
)

// Engine used to power all defined renderers.
// This allows you to configure the system to your
// prefered settings, instead of just getting
// the defaults.
type Engine struct {
	Options
}

// New render.Engine ready to go with your Options
// and some defaults we think you might like.
func New(opts Options) *Engine {
	e := &Engine{
		Options: opts,
	}
	e.RegisterHelpers(helpers.Helpers)
	return e
}

// RegisterHelper adds a helper to a template with the given name.
// See github.com/aymerick/raymond for more details on helpers.
/*
	e.RegisterHelper("upcase", strings.ToUpper)
*/
func (e *Engine) RegisterHelper(name string, helper interface{}) {
	e.RegisterHelpers(map[string]interface{}{
		name: helper,
	})
}

// RegisterHelpers adds helpers to a template with the given name.
// See github.com/aymerick/raymond for more details on helpers.
/*
	h := map[string]interface{}{
		"upcase": strings.ToUpper,
		"downcase": strings.ToLower,
	}
	e.RegisterHelpers(h)
*/
func (e *Engine) RegisterHelpers(helpers map[string]interface{}) {
	defer func() {
		// Since raymond panics(!!) when a helper is already registered
		// let's recover and move on.
		recover()
	}()
	raymond.RegisterHelpers(helpers)
}
