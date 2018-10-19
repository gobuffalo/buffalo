package meta

import (
	"github.com/gobuffalo/meta"
	"github.com/markbates/oncer"
)

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo/meta", "Use github.com/gobuffalo/meta instead.")
}

// App represents meta data for a Buffalo application on disk
// Use meta.App instead
type App = meta.App

// New App based on the details found at the provided root path
// Use meta.New intead.
var New = meta.New

// ResolveSymlinks takes a path and gets the pointed path
// if the original one is a symlink.
// Use meta.ResolveSymlinks instead
var ResolveSymlinks = meta.ResolveSymlinks
