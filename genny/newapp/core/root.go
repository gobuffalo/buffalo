package core

import (
	"go/build"
	"html/template"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func rootGenerator(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	// validate opts
	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	g.Transformer(genny.Dot())

	// add common templates
	if err := g.Box(packr.New("buffalo:genny:newapp:core", "../core/templates")); err != nil {
		return g, errors.WithStack(err)
	}

	data := map[string]interface{}{
		"opts": opts,
	}

	helpers := template.FuncMap{}

	t := gotools.TemplateTransformer(data, helpers)
	g.Transformer(t)

	if !opts.App.WithModules {
		c := build.Default
		g.RunFn(validateInGoPath(c.SrcDirs()))
	}

	g.RunFn(func(r *genny.Runner) error {
		f := genny.NewFile("config/buffalo-app.toml", nil)
		if err := toml.NewEncoder(f).Encode(opts.App); err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	})

	return g, nil
}
