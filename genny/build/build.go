package build

import (
	"path/filepath"
	"strings"
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

	box := packr.NewBox("../build/templates")
	if err := g.Box(box); err != nil {
		return g, errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	ctx.Set("buildTime", opts.BuildTime.Format(time.RFC3339))
	ctx.Set("buildVersion", opts.BuildVersion)

	// create the ./a pkg
	g.RunFn(apkg(opts))
	g.Transformer(plushgen.Transformer(ctx))

	// clean up everything!
	g.RunFn(func(r *genny.Runner) error {
		var err error
		opts.rollback.Range(func(k, v interface{}) bool {
			f := genny.NewFile(k.(string), strings.NewReader(v.(string)))
			r.Logger.Debug("rolling back modified file", f.Name())
			if err = r.File(f); err != nil {
				return false
			}
			r.Disk.Remove(f.Name())
			return true
		})
		if err != nil {
			return errors.WithStack(err)
		}
		for _, f := range r.Disk.Files() {
			r.Logger.Debug("cleaning up generated file", f.Name())
			if err := r.Disk.Delete(f.Name()); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	return g, nil
}
