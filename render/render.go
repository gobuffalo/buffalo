package render

import "sync"

// Engine used to power all defined renderers.
// This allows you to configure the system to your
// prefered settings, instead of just getting
// the defaults.
type Engine struct {
	Options
	moot *sync.Mutex
}

// New render.Engine ready to go with your Options
// and some defaults we think you might like.
func New(opts Options) *Engine {
	opts.TemplateHelpers = map[string]interface{}{}

	e := &Engine{
		Options: opts,
		moot:    &sync.Mutex{},
	}
	return e
}
