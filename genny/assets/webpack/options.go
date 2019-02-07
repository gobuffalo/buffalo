package webpack

import (
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
)

// Options for creating a new webpack setup
type Options struct {
	meta.App
	Bootstrap int `json:"bootstrap"`
}

// Validate options
func (opts *Options) Validate() error {
	if opts.Bootstrap == 0 {
		opts.Bootstrap = 4
	}
	bs := opts.Bootstrap
	if bs < 3 && bs > 4 {
		return errors.Errorf("unknown bootstrap version %d", bs)
	}
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	return nil
}
