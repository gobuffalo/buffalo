package build

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/genny/assets/webpack"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"

	"github.com/gobuffalo/packr/v2/jam"
	"github.com/gobuffalo/packr/v2/jam/parser"
)

func assets(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if opts.App.WithNodeJs || opts.App.WithWebpack {
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
		g.RunFn(func(r *genny.Runner) error {
			tool := "yarnpkg"
			if !opts.App.WithYarn {
				tool = "npm"
			}

			c := exec.CommandContext(r.Context, tool, "run", "build")
			if _, err := opts.App.NodeScript("build"); err != nil {
				// Fallback on legacy runner
				c = exec.CommandContext(r.Context, webpack.BinPath)
			}

			bb := &bytes.Buffer{}
			c.Stdout = bb
			c.Stderr = bb

			if err := r.Exec(c); err != nil {
				r.Logger.Error(bb.String())
				return err
			}
			return nil

		})
	}

	g.RunFn(func(r *genny.Runner) error {
		ro := &parser.RootsOptions{}

		if !opts.WithAssets {
			ro.Ignores = append(ro.Ignores, "public/assets")
		}

		opts := jam.PackOptions{
			Roots:        []string{opts.App.Root},
			RootsOptions: ro,
		}
		return jam.Pack(opts)
	})

	if opts.ExtractAssets && opts.WithAssets {
		// mount the archived assets generator
		aa, err := archivedAssets(opts)
		if err != nil {
			return g, err
		}
		g.Merge(aa)
	}

	return g, nil
}
