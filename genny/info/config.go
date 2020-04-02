package info

import (
	"path"

	"github.com/gobuffalo/genny/v2"
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
		return box.Walk(func(p string, f packr.File) error {
			opts.Out.Header("Buffalo: " + path.Join("config", p))
			opts.Out.WriteString(f.String() + "\n")
			return nil
		})
	}
}
