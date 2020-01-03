package api

import (
	"html/template"

	"github.com/gobuffalo/buffalo/genny/newapp/core"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen"
	"github.com/gobuffalo/genny/gogen/gomods"
	"github.com/gobuffalo/packr/v2"
)

// New generator for creating a Buffalo API application
func New(opts *Options) (*genny.Group, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	gg, err := core.New(opts.Options)
	if err != nil {
		return gg, err
	}

	g := genny.New()
	data := map[string]interface{}{
		"opts": opts,
	}

	helpers := template.FuncMap{}

	t := gogen.TemplateTransformer(data, helpers)
	g.Transformer(t)
	g.Box(packr.New("buffalo:genny:newapp:api", "../api/templates"))

	gg.Add(g)

	g, err := gomods.Tidy(opts.App.Root, false)
	if err != nil {
		return gg, err
	}
	gg.Add(g)

	return gg, nil
}
