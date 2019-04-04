package webpack

import (
	"github.com/gobuffalo/meta"
)

// Options for creating a new webpack setup
type Options struct {
	meta.App
}

// Validate options
func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	return nil
}
