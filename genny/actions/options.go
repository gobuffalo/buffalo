package actions

import (
	"fmt"

	"github.com/gobuffalo/meta"
)

// Options for the actions generator
type Options struct {
	App           meta.App
	Name          string
	Actions       []string
	Method        string
	SkipTemplates bool
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return fmt.Errorf("you must provide a name")
	}

	if len(opts.Actions) == 0 {
		return fmt.Errorf("you must provide at least one action name")
	}

	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}

	if len(opts.Method) == 0 {
		opts.Method = "GET"
	}
	return nil
}
