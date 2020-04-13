package plugins

import (
	"bytes"

	"github.com/gobuffalo/buffalo/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
)

// NewEncodePluginsRunner will return a runner that will encode the plugins file
func NewEncodePluginsRunner(app meta.App, plugs *plugdeps.Plugins) func(r *genny.Runner) error {
	return func(r *genny.Runner) error {
		p := plugdeps.ConfigPath(app)
		bb := &bytes.Buffer{}
		if err := plugs.Encode(bb); err != nil {
			return err
		}

		return r.File(genny.NewFile(p, bb))
	}
}
