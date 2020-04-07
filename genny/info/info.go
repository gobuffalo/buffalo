package info

import (
	"path/filepath"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/packr/v2"
)

// New returns a generator that performs buffalo
// related rx checks
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.RunFn(appDetails(opts))

	cBox := packr.Folder(filepath.Join(opts.App.Root, "config"))
	g.RunFn(configs(opts, cBox))

	aBox := packr.Folder(opts.App.Root)
	g.RunFn(pkgChecks(opts, aBox))

	return g, nil
}
