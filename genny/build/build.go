package build

import (
	"path/filepath"
	"time"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}
	g.Transformer(genny.Dot())

	// validate templates
	tb := packr.NewBox(filepath.Join(opts.App.Root, "templates"))
	g.RunFn(ValidateTemplates(tb, opts.TemplateValidators))

	// rename main() to originalMain()
	g.RunFn(transformMain(opts))

	// add any necessary templates for the build
	box := packr.NewBox("../build/templates")
	if err := g.Box(box); err != nil {
		return g, errors.WithStack(err)
	}

	// configure plush
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	ctx.Set("buildTime", opts.BuildTime.Format(time.RFC3339))
	ctx.Set("buildVersion", opts.BuildVersion)
	g.Transformer(plushgen.Transformer(ctx))

	// create the ./a pkg
	g.RunFn(apkg(opts))

	// clean up everything!
	g.RunFn(cleanup(opts))

	return g, nil
}
