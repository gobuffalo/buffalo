package info

import (
	"os"

	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/meta"
)

type Options struct {
	App meta.App
	Out rx.Writer
	// add your stuff here
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