package resource

import (
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.New("github.com/gobuffalo/buffalo/genny/resource/templates", "../resource/templates")); err != nil {
		return g, errors.WithStack(err)
	}

	// data := gotools.Data{}
	data := map[string]interface{}{
		"opts": opts,
	}
	helpers := template.FuncMap{}
	g.Transformer(gotools.TemplateTransformer(data, helpers))

	return g, nil
}
