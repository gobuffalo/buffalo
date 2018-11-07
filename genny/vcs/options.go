package vcs

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/meta"
)

// Available VCS implementations
var Available = []string{"git", "bzr", "none"}

// Options for VCS generator
type Options struct {
	App      meta.App
	Provider string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}

	var found bool
	for _, a := range Available {
		if opts.Provider == a {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("unknown provider %q expecting one of %s", opts.Provider, strings.Join(Available, ", "))
	}
	return nil
}
