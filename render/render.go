package render

import (
	"github.com/gobuffalo/helpers"
	"github.com/gobuffalo/helpers/forms"
	"github.com/gobuffalo/helpers/forms/bootstrap"
	"github.com/gobuffalo/plush/v4"
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
		opts.Helpers = defaultHelpers()
	}

	if opts.TemplateEngines == nil {
		opts.TemplateEngines = map[string]TemplateEngine{}
	}
	if _, ok := opts.TemplateEngines["html"]; !ok {
		opts.TemplateEngines["html"] = plush.BuffaloRenderer
	}
	if _, ok := opts.TemplateEngines["plush"]; !ok {
		opts.TemplateEngines["plush"] = plush.BuffaloRenderer
	}
	if _, ok := opts.TemplateEngines["text"]; !ok {
		opts.TemplateEngines["text"] = plush.BuffaloRenderer
	}
	if _, ok := opts.TemplateEngines["txt"]; !ok {
		opts.TemplateEngines["txt"] = plush.BuffaloRenderer
	}
	if _, ok := opts.TemplateEngines["js"]; !ok {
		opts.TemplateEngines["js"] = plush.BuffaloRenderer
	}
	if _, ok := opts.TemplateEngines["md"]; !ok {
		opts.TemplateEngines["md"] = MDTemplateEngine
	}
	if _, ok := opts.TemplateEngines["tmpl"]; !ok {
		opts.TemplateEngines["tmpl"] = GoTemplateEngine
	}

	if opts.DefaultContentType == "" {
		opts.DefaultContentType = "text/html; charset=utf-8"
	}

	e := &Engine{
		Options: opts,
	}
	return e
}

func defaultHelpers() Helpers {
	h := Helpers(helpers.ALL())
	h[forms.FormKey] = bootstrap.Form
	h[forms.FormForKey] = bootstrap.FormFor
	h["form_for"] = bootstrap.FormFor
	return h
}
