package actions

import (
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
)

type Options struct {
	App          meta.App
	Name         string
	Method       string
	SkipTemplate bool
	Actions      []string
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if len(opts.Name) == 0 {
		return errors.New("you must provide a name")
	}

	if len(opts.Actions) == 0 {
		return errors.New("you must provide at least one action name")
	}

	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}

	if len(opts.Method) == 0 {
		opts.Method = "GET"
	}
	return nil
}
