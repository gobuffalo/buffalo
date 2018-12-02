package build

import (
	"context"
	"io/ioutil"
	"os/exec"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	pack "github.com/gobuffalo/packr/builder"
	"github.com/pkg/errors"
)

func assets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if opts.App.WithNodeJs {
		if _, err := opts.App.NodeScript("build"); err != nil {
			return nil, err
		}
		g.RunFn(func(r *genny.Runner) error {
			r.Logger.Debugf("setting NODE_ENV = %s", opts.Environment)
			return envy.MustSet("NODE_ENV", opts.Environment)
		})
		tool := "yarnpkg"
		if !opts.App.WithYarn {
			tool = "npm"
		}
		if _, err := exec.LookPath(tool); err != nil {
			return nil, errors.Errorf("couldn't find %s tool", tool)
		}
		cmd := exec.Command(tool, "run", "build")
		cmd.Stdout = ioutil.Discard
		g.Command(cmd)
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
			return g, errors.WithStack(err)
		}
		g.Merge(aa)
	}

	g.RunFn(func(r *genny.Runner) error {
		return p.Run()
	})

	return g, nil
}
