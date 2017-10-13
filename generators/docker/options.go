package docker

import (
	"github.com/gobuffalo/buffalo/meta"
)

// Options for generating a new docker file
type Options struct {
	App     meta.App
	Version string
	Style   string
	AsWeb   bool
}

// NewOptions returns a well formed set of options for generating a docker file
func NewOptions() Options {
	o := Options{
		App:     meta.New("."),
		Version: "latest",
		Style:   "multi",
	}
	o.AsWeb = o.App.WithWebpack

	return o
}
