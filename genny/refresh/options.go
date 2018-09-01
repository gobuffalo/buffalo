package refresh

import "github.com/gobuffalo/buffalo/meta"

// Options for creating a new refresh config
type Options struct {
	App meta.App
}

// Validate options
func (opts *Options) Validate() error {
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}
	return nil
}
