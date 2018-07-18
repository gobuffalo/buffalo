package docker

import (
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/buffalo/runtime"
)

// Generator for generating a new docker file
type Generator struct {
	App           meta.App `json:"app"`
	Version       string   `json:"version"`
	Style         string   `json:"style"`
	AsWeb         bool     `json:"as_web"`
	DBType        string   `json:"db_type"`
	DockerCompose string   `json:"docker_compose"`
}

// New returns a well formed set of options for generating a docker file
func New() Generator {
	o := Generator{
		App:     meta.New("."),
		Version: runtime.Version,
		Style:   "multi",
		DBType:  "postgres",
	}
	o.AsWeb = o.App.WithWebpack

	return o
}
