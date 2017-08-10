package render

import (
	"github.com/gobuffalo/plush"
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

	if opts.TemplateEngine == nil {
		opts.TemplateEngine = plush.BuffaloRenderer
	}

	e := &Engine{
		Options: opts,
	}
	return e
}
