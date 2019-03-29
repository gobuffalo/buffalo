package info

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2"
)

func configs(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		box := packr.Folder(filepath.Join(opts.App.Root, "config"))
		return box.Walk(func(path string, f packr.File) error {
			p := strings.TrimPrefix(path, opts.App.Root)
			p = strings.TrimPrefix(p, string(filepath.Separator))
			opts.Out.WriteString(fmt.Sprintf("\n### %s\n", p))
			opts.Out.WriteString(f.String())
			return nil
		})
	}
}
