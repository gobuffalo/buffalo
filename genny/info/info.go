package info

import (
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	g.RunFn(appDetails(opts))
	box := packr.Folder(filepath.Join(opts.App.Root, "config"))
	g.RunFn(configs(opts, box))

	return g, nil
}
