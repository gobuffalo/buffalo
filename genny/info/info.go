package info

import (
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	g.RunFn(appDetails(opts))
	g.RunFn(configs(opts))

	return g, nil
}
