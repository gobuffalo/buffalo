package grift

import (
	"github.com/markbates/inflect"
)

// Options for creating a new grift task
type Options struct {
	Name       inflect.Name   `json:"name"`
	Parts      []inflect.Name `json:"parts"`
	Args       []string       `json:"args"`
	Namespaced bool           `json:"namespaced"`
}

// Last checks if the name is the last of the parts
func (g Options) Last(n inflect.Name) bool {
	return g.Parts[len(g.Parts)-1] == n
}
