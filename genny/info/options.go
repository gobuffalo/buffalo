package info

import (
	"io"
	"os"

	"github.com/gobuffalo/meta"
)

type Options struct {
	App meta.App
	Out io.Writer
	// add your stuff here
}

// Validate that options are usuable
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	return nil
}
