package webpack

import (
	"fmt"

	"github.com/gobuffalo/meta"
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
		return fmt.Errorf("unknown bootstrap version %d", bs)
	}
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	return nil
}
