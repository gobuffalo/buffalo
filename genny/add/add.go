package add

import (
	"bytes"
	"path/filepath"

	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	bb := &bytes.Buffer{}
	plugs := plugdeps.New()
	plugs.Add(opts.Plugins...)
	if err := plugs.Encode(bb); err != nil {
		return g, errors.WithStack(err)
	}

	cpath := filepath.Join(opts.App.Root, "config", "buffalo-plugins.toml")
	g.File(genny.NewFile(cpath, bb))

	return g, nil
}
