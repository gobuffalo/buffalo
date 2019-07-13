package install

import (
	"os"

	"github.com/gobuffalo/buffalo/internal/takeon/github.com/markbates/errx"
	"github.com/gobuffalo/buffalo/plugins/plugdeps"
	"github.com/gobuffalo/meta"
)

// Options container for passing needed info for
// installing plugins and adding them to the config file.
type Options struct {
	App     meta.App
	Plugins []plugdeps.Plugin
	Tags    meta.BuildTags
	Vendor  bool
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		opts.App = meta.New(pwd)
	}
	if len(opts.Plugins) == 0 {
		plugs, err := plugdeps.List(opts.App)
		if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
			return err
		}
		opts.Plugins = plugs.List()
	}

	for i, p := range opts.Plugins {
		p.Tags = opts.App.BuildTags("", append(opts.Tags, p.Tags...)...)
		opts.Plugins[i] = p
	}

	return nil
}
