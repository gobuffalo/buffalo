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

	if opts.TemplateEngines == nil {
		opts.TemplateEngines = map[string]TemplateEngine{}
	}
	if _, ok := opts.TemplateEngines["html"]; !ok {
		opts.TemplateEngines["html"] = plush.BuffaloRenderer
	}
	if _, ok := opts.TemplateEngines["text"]; !ok {
		opts.TemplateEngines["text"] = plush.BuffaloRenderer
	}
	if _, ok := opts.TemplateEngines["js"]; !ok {
		opts.TemplateEngines["js"] = JSTemplateEngine
	}
	if _, ok := opts.TemplateEngines["md"]; !ok {
		opts.TemplateEngines["md"] = MDTemplateEngine
	}
	if _, ok := opts.TemplateEngines["tmpl"]; !ok {
		opts.TemplateEngines["tmpl"] = GoTemplateEngine
	}

	e := &Engine{
		Options: opts,
	}
	return e
}
