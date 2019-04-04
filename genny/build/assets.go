package build

import (
	"context"
	"io/ioutil"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/genny/assets/webpack"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"

	pack "github.com/gobuffalo/packr/builder"
)

func assets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if opts.App.WithWebpack {
		if opts.CleanAssets {
			g.RunFn(func(r *genny.Runner) error {
				r.Delete(filepath.Join(opts.App.Root, "public", "assets"))
				return nil
			})
		}
		g.RunFn(func(r *genny.Runner) error {
			r.Logger.Debugf("setting NODE_ENV = %s", opts.Environment)
			return envy.MustSet("NODE_ENV", opts.Environment)
		})
		c := exec.Command(webpack.BinPath)
		c.Stdout = ioutil.Discard
		g.Command(c)
	}

	p := pack.New(context.Background(), opts.App.Root)
	p.Compress = true

	if !opts.WithAssets {
		p.IgnoredBoxes = append(p.IgnoredBoxes, "../public/assets")
	} else {
		p.IgnoredFolders = p.IgnoredFolders[1:]
	}

	if opts.ExtractAssets && opts.WithAssets {
		p.IgnoredBoxes = append(p.IgnoredBoxes, "../public/assets")
		// mount the archived assets generator
		aa, err := archivedAssets(opts)
		if err != nil {
			return g, err
		}
		g.Merge(aa)
	}

	g.RunFn(func(r *genny.Runner) error {
		return p.Run()
	})

	return g, nil
}
