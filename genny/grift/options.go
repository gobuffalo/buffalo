package grift

import (
	"strings"

	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

// Options for creating a new grift task
type Options struct {
	Name       inflect.Name   `json:"name"`
	Parts      []inflect.Name `json:"parts"`
	Args       []string       `json:"args"`
	Namespaced bool           `json:"namespaced"`
}

// Last checks if the name is the last of the parts
func (opts Options) Last(n inflect.Name) bool {
	return opts.Parts[len(opts.Parts)-1] == n
}

// Validate options
func (opts *Options) Validate() error {
	if len(opts.Args) == 0 {
		return errors.New("you need to provide a name for the grift task")
	}

	opts.Namespaced = strings.Contains(opts.Args[0], ":")

	for _, n := range strings.Split(opts.Args[0], ":") {
		opts.Parts = append(opts.Parts, inflect.Name(n))
	}
	opts.Name = opts.Parts[len(opts.Parts)-1]
	return nil
}
