package grift

import (
	"strings"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

// Generator for creating a new grift task
type Generator struct {
	App        meta.App       `json:"app"`
	Name       inflect.Name   `json:"name"`
	Parts      []inflect.Name `json:"parts"`
	Args       []string       `json:"args"`
	Namespaced bool           `json:"namespaced"`
}

// Last checks if the name is the last of the parts
func (g Generator) Last(n inflect.Name) bool {
	return g.Parts[len(g.Parts)-1] == n
}

// New generator for grift tasks
func New(args ...string) (Generator, error) {
	g := Generator{
		App:   meta.New("."),
		Args:  args,
		Parts: []inflect.Name{},
	}
	if len(args) > 0 {
		g.Namespaced = strings.Contains(args[0], ":")

		for _, n := range strings.Split(args[0], ":") {
			g.Parts = append(g.Parts, inflect.Name(n))
		}
		g.Name = inflect.Name(g.Parts[len(g.Parts)-1])
	}

	return g, g.Validate()
}

// Validate the generator
func (g Generator) Validate() error {
	if len(g.Args) < 1 {
		return errors.New("you need to provide a name for the grift task")
	}
	return nil
}
