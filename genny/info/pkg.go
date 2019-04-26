package info

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packd"
)

func pkgChecks(opts *Options, box packd.Finder) genny.RunFn {
	return func(r *genny.Runner) error {
		for _, x := range []string{"go.mod", "Gopkg.toml", "Gopkg.lock"} {
			f, err := box.FindString(x)
			if err == nil {
				opts.Out.Header("\nBuffalo: " + x)
				opts.Out.WriteString(f)
			}
		}
		return nil
	}
}
