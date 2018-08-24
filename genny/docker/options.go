package docker

import (
	"github.com/gobuffalo/buffalo/meta"
)

// Options for generating a new docker file
type Options struct {
	App     meta.App `json:"app"`
	Version string   `json:"version"`
	Style   string   `json:"style"`
	AsWeb   bool     `json:"as_web"`
}
