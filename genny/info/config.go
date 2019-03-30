package info

import (
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
)

// ListWalker allows for a box that supports listing and walking
type ListWalker interface {
	packd.Lister
	packd.Walkable
}

func configs(opts *Options, box ListWalker) genny.RunFn {
	return func(r *genny.Runner) error {
		if len(box.List()) == 0 {
			return nil
		}
		return box.Walk(func(path string, f packr.File) error {
			opts.Out.Header("Buffalo: " + filepath.Join("config", path))
			opts.Out.WriteString(f.String() + "\n")
			return nil
		})
	}
}
