package core

import (
	"go/build"
	"html/template"
	"os/exec"

	"github.com/BurntSushi/toml"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
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
	if err := g.Box(packr.NewBox("../core/templates")); err != nil {
		return g, errors.WithStack(err)
	}

	if opts.App.WithModules {
		// add go.mod templates
		g.Command(exec.Command(genny.GoBin(), "mod", "init", opts.App.PackagePkg))
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
