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

// See github.com/aymerick/raymond for more details on helpers.
func (e *Engine) RegisterHelper(name string, helper interface{}) {
	raymond.RegisterHelper(name, helper)
}

// See github.com/aymerick/raymond for more details on helpers.
func (e *Engine) RegisterHelpers(helpers map[string]interface{}) {
	defer func() {
		// Since raymond panics(!!) when a helper is already registered
		// let's recover and move on.
		recover()
	}()
	raymond.RegisterHelpers(helpers)
}
