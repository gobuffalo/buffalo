package build

import (
	"time"

	"github.com/gobuffalo/events"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/pkg/errors"
)

// New generator for building a Buffalo application
// This powers the `buffalo build` command and can be
// used to programatically build/compile Buffalo
// applications.
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}
	g.ErrorFn = func(err error) {
		events.EmitError(EvtBuildStopErr, err, events.Payload{"opts": opts})
	}

	g.Event(EvtBuildStart, events.Payload{"opts": opts})

	g.Transformer(genny.Dot())

	// validate templates
	g.RunFn(ValidateTemplates(templateWalker(opts.App), opts.TemplateValidators))

	// rename main() to originalMain()
	g.RunFn(transformMain(opts))

	// add any necessary templates for the build
	box := packr.New("buffalo:genny:build", "../build/templates")
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
	ag, err := apkg(opts)
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.Merge(ag)

	if opts.WithAssets {
		// mount the assets generator
		ag, err := assets(opts)
		if err != nil {
			return g, errors.WithStack(err)
		}
		g.Merge(ag)
	}

	// mount the build time dependency generator
	dg, err := buildDeps(opts)
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.Merge(dg)

	// create the final go build command
	c, err := buildCmd(opts)
	if err != nil {
		return g, errors.WithStack(err)
	}
	g.Command(c)

	// clean up everything!
	g.RunFn(cleanup(opts))

	g.Event(EvtBuildStop, events.Payload{"opts": opts})
	return g, nil
}
