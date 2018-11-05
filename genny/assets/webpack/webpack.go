package webpack

import (
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/buffalo/genny/assets/standard"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BinPath is the path to the local install of webpack
var BinPath = filepath.Join("node_modules", ".bin", "webpack")

// New generator for create webpack asset files
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if _, err := exec.LookPath("npm"); err != nil {
		logrus.Info("Could not find npm. Skipping webpack generation.")
		return standard.New(&standard.Options{})
	}

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	g.Box(packr.NewBox("../webpack/templates"))

	data := map[string]interface{}{
		"opts": opts,
	}
	t := gotools.TemplateTransformer(data, gotools.TemplateHelpers)
	g.Transformer(t)
	g.Transformer(genny.Dot())

	g.RunFn(func(r *genny.Runner) error {
		return installPkgs(r, opts)
	})

	return g, nil
}

func installPkgs(r *genny.Runner, opts *Options) error {
	command := "yarnpkg"

	if !opts.App.WithYarn {
		command = "npm"
	} else {
		if err := installYarn(r); err != nil {
			return errors.WithStack(err)
		}
	}
	args := []string{"install", "--no-progress", "--save"}

	c := exec.Command(command, args...)
	c.Stdout = yarnWriter{
		fn: r.Logger.Debug,
	}
	c.Stderr = yarnWriter{
		fn: r.Logger.Debug,
	}
	return r.Exec(c)
}

type yarnWriter struct {
	fn func(...interface{})
}

func (y yarnWriter) Write(p []byte) (int, error) {
	y.fn(string(p))
	return len(p), nil
}

func installYarn(r *genny.Runner) error {
	// if there's no yarn, install it!
	if _, err := exec.LookPath("yarnpkg"); err == nil {
		return nil
	}
	yargs := []string{"install", "-g", "yarn"}
	return r.Exec(exec.Command("npm", yargs...))
}
