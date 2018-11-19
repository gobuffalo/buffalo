package build

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo/generators/assets/webpack"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"

	pack "github.com/gobuffalo/packr/builder"
)

func assets(opts *Options) (*genny.Generator, error) {

	g := genny.New()

	// Quick way to ensure this directory exists so static content is created.
	if opts.ExtractAssets {
		g.File(genny.NewDir("public/assets", os.FileMode(int(0777))))
	}
	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}
	if opts.App.WithWebpack {
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
		// TODO: note the difference between ../public/assets and public/assets
		p.IgnoredBoxes = append(p.IgnoredBoxes, "../public/assets")
	} else {
		p.IgnoredFolders = p.IgnoredFolders[1:]
	}

	if opts.ExtractAssets && opts.WithAssets && !opts.AsAPI {
		p.IgnoredBoxes = append(p.IgnoredBoxes, "../public/assets")
		// mount the archived assets generator
		aa, err := archivedAssets(opts)
		if err != nil {
			return g, errors.WithStack(err)
		}
		g.Merge(aa)
	}
	g.RunFn(func(r *genny.Runner) error {
		return p.Run()
	})
	return g, nil
}
