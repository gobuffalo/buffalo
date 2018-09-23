package plugins

import (
	"github.com/gobuffalo/buffalo-plugins/plugins"
	"github.com/markbates/oncer"
)

// List is deprecated and moved to github.com/gobuffalo/buffalo-plugins/plugins
type List plugins.List

// Available is deprecated and moved to github.com/gobuffalo/buffalo-plugins/plugins
var Available = plugins.Available

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo/plugins", "Use github.com/gobuffalo/buffalo-plugsin/plugins instead.")
}
