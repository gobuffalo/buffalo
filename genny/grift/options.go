package grift

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/flect/name"
)

// Options for creating a new grift task
type Options struct {
	Name       name.Ident   `json:"name"`
	Parts      []name.Ident `json:"parts"`
	Args       []string     `json:"args"`
	Namespaced bool         `json:"namespaced"`
}

// Last checks if the name is the last of the parts
func (opts Options) Last(n name.Ident) bool {
	return opts.Parts[len(opts.Parts)-1].String() == n.String()
}

// Validate options
func (opts *Options) Validate() error {
	if len(opts.Args) == 0 {
		return fmt.Errorf("you need to provide a name for the grift task")
	}

	opts.Namespaced = strings.Contains(opts.Args[0], ":")

	for _, n := range strings.Split(opts.Args[0], ":") {
		opts.Parts = append(opts.Parts, name.New(n))
	}
	opts.Name = opts.Parts[len(opts.Parts)-1]
	return nil
}
