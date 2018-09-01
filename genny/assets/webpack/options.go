package webpack

import (
	"github.com/gobuffalo/buffalo/meta"
	"github.com/pkg/errors"
)

// Options for creating a new webpack setup
type Options struct {
	meta.App
	Bootstrap int `json:"bootstrap"`
}

func (opts *Options) Validate() error {
	if opts.Bootstrap == 0 {
		opts.Bootstrap = 4
	}
	bs := opts.Bootstrap
	if bs < 3 && bs > 4 {
		return errors.Errorf("unknown bootstrap version %d", bs)
	}
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}
	return nil
}
