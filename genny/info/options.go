package info

import (
	"os"

	"github.com/gobuffalo/clara/v2/genny/rx"
	"github.com/gobuffalo/meta"
)

// Options for the info generator
type Options struct {
	App meta.App
	Out rx.Writer
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	if opts.Out.Writer == nil {
		opts.Out = rx.NewWriter(os.Stdout)
	}
	return nil
}
