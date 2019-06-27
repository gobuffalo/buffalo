package add

import (
	"bytes"
	"path/filepath"

	"github.com/gobuffalo/buffalo/plugins/plugdeps"
	"github.com/gobuffalo/genny"
)

// New add plugin to the config file
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	bb := &bytes.Buffer{}
	plugs := plugdeps.New()
	plugs.Add(opts.Plugins...)
	if err := plugs.Encode(bb); err != nil {
		return g, err
	}

	cpath := filepath.Join(opts.App.Root, "config", "buffalo-plugins.toml")
	g.File(genny.NewFile(cpath, bb))

	return g, nil
}
