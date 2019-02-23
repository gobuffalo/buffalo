package web

import (
	"html/template"

	"github.com/gobuffalo/buffalo/genny/assets/standard"
	"github.com/gobuffalo/buffalo/genny/assets/webpack"
	"github.com/gobuffalo/buffalo/genny/newapp/core"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

// New generator for creating a Buffalo Web application
func New(opts *Options) (*genny.Group, error) {
	if err := opts.Validate(); err != nil {
		return nil, errors.WithStack(err)
	}

	gg, err := core.New(opts.Options)
	if err != nil {
		return gg, errors.WithStack(err)
	}

	g := genny.New()
	g.Transformer(genny.Dot())
	data := map[string]interface{}{
		"opts": opts,
	}

	helpers := template.FuncMap{}

	t := gogen.TemplateTransformer(data, helpers)
	g.Transformer(t)
	g.Box(packr.New("buffalo:genny:newapp:web", "../web/templates"))

	gg.Add(g)

	if opts.Webpack != nil {
		// add the webpack generator
		g, err = webpack.New(opts.Webpack)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(g)
	}

	if opts.Standard != nil {
		// add the standard generator
		g, err = standard.New(opts.Standard)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(g)
	}

	return gg, nil
}
