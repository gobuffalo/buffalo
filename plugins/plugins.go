package plugins

import (
	"fmt"

	"github.com/gobuffalo/buffalo-plugins/plugins"
)

// List is deprecated and moved to github.com/gobuffalo/buffalo-plugins/plugins
type List plugins.List

// Available is deprecated and moved to github.com/gobuffalo/buffalo-plugins/plugins
var Available = plugins.Available

func init() {
	fmt.Println("github.com/gobuffalo/buffalo/plugins has been deprecated in v0.12.4, and will be removed in v0.13.0. Use github.com/gobuffalo/buffalo-plugins/plugins directly.")
}
