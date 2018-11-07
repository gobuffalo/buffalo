package web

import (
	"github.com/gobuffalo/buffalo/genny/assets/standard"
	"github.com/gobuffalo/buffalo/genny/assets/webpack"
	"github.com/gobuffalo/buffalo/genny/newapp/core"
	"github.com/pkg/errors"
)

// Options for a web app
type Options struct {
	*core.Options
	Webpack  *webpack.Options
	Standard *standard.Options
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.Options == nil {
		opts.Options = &core.Options{}
	}

	if err := opts.Options.Validate(); err != nil {
		return errors.WithStack(err)
	}

	if opts.Docker != nil {
		if opts.Docker.App.IsZero() {
			opts.Docker.App = opts.App
		}
		if err := opts.Docker.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.Webpack != nil {
		if opts.Webpack.App.IsZero() {
			opts.Webpack.App = opts.App
		}
		if err := opts.Webpack.Validate(); err != nil {
			return errors.WithStack(err)
		}
	}

	if opts.Standard != nil && opts.Webpack != nil {
		return errors.New("you can not use both webpack and standard generators")
	}

	return nil
}
