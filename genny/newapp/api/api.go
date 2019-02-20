package api

import (
	"html/template"

	"github.com/gobuffalo/buffalo/genny/newapp/core"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

// New generator for creating a Buffalo API application
func New(opts *Options) (*genny.Group, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	gg, err := core.New(opts.Options)
	if err != nil {
		return gg, errors.WithStack(err)
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

	return gg, nil
}
