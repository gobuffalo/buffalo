package resource

import (
	"os/exec"
	"text/template"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	core := packr.New("github.com/gobuffalo/buffalo/genny/resource/templates/core", "../resource/templates/core")
	if err := g.Box(core); err != nil {
		return g, errors.WithStack(err)
	}

	var abox packd.Box
	if opts.UseModel {
		abox = packr.New("github.com/gobuffalo/buffalo/genny/resource/templates/use_model", "../resource/templates/use_model")
	} else {
		abox = packr.New("github.com/gobuffalo/buffalo/genny/resource/templates/standard", "../resource/templates/standard")
	}

	if err := g.Box(abox); err != nil {
		return g, errors.WithStack(err)
	}
	pres := presenter{
		App:   opts.App,
		Name:  name.New(opts.Name),
		Model: name.New(opts.Model),
		Attrs: opts.Attrs,
	}
	data := map[string]interface{}{
		"opts": pres,
	}
	helpers := template.FuncMap{
		"camelize": func(s string) string {
			return flect.Camelize(s)
		},
	}

	x := pres.Name.File().String()
	g.Transformer(gotools.TemplateTransformer(data, helpers))
	g.Transformer(genny.Replace("resource-name", x))
	g.Transformer(genny.Replace("resource-use_model", x))

	g.RunFn(func(r *genny.Runner) error {
		if !opts.SkipModel && !opts.UseModel {
			if _, err := r.LookPath("buffalo-pop"); err != nil {
				if err := gotools.Get("github.com/gobuffalo/buffalo-pop")(r); err != nil {
					return errors.WithStack(err)
				}
			}

			return r.Exec(modelCommand(pres.Model, opts))
		}

		return nil
	})
	return g, nil
}

func modelCommand(model name.Ident, opts *Options) *exec.Cmd {
	args := opts.Attrs.Slice()
	args = append(args[:0], args[0+1:]...)
	args = append([]string{"pop", "g", "model", model.Singularize().Underscore().String()}, args...)

	if opts.SkipMigration {
		args = append(args, "--skip-migration")
	}

	return exec.Command("buffalo-pop", args...)
}

type presenter struct {
	App   meta.App
	Name  name.Ident
	Model name.Ident
	Attrs attrs.Attrs
	// SkipMigration bool
	// SkipModel     bool
	// SkipTemplates bool
	// UseModel      bool
	// Args          []string
}
