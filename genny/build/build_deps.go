package build

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/dep"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
)

func buildDeps(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if opts.App.WithDep {
		// mount the dep generator
		dg, err := dep.Ensure(false)
		if err != nil {
			return g, errors.WithStack(err)
		}
		g.Merge(dg)
	} else {
		// mount the go get runner
		g.RunFn(gotools.Get("./..."))
	}
	return g, nil
}
