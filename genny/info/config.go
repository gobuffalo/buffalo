package info

import (
	"fmt"
	"path/filepath"
	"strings"

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
		opts.Out.Header("Buffalo: Config Files")
		return box.Walk(func(path string, f packr.File) error {
			p := strings.TrimPrefix(path, opts.App.Root)
			p = strings.TrimPrefix(p, string(filepath.Separator))
			opts.Out.WriteString(fmt.Sprintf("\n### %s\n", p))
			opts.Out.WriteString(f.String())
			return nil
		})
	}
}
