package webpack

import (
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/gobuffalo/buffalo/genny/assets/standard"
	"github.com/gobuffalo/buffalo/meta"
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

	if opts.Bootstrap == 0 {
		opts.Bootstrap = 4
	}
	bs := opts.Bootstrap
	if bs < 3 && bs > 4 {
		return nil, errors.Errorf("bootstrap can only be 3 or 4 not %d", bs)
	}
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}

	g.Box(packr.NewBox("../webpack/templates"))
	data := map[string]interface{}{
		"opts": opts,
	}

	g.RunFn(func(r *genny.Runner) error {
		command := "yarnpkg"

		if !opts.App.WithYarn {
			command = "npm"
		} else {
			if err := installYarn(r); err != nil {
				return errors.WithStack(err)
			}
		}
		args := []string{"install", "--no-progress", "--save"}
		return r.Exec(exec.Command(command, args...))
	})

	h := template.FuncMap{}
	t := gotools.TemplateTransformer(data, h)
	g.Transformer(t)
	g.Transformer(genny.Dot())
	return g, nil
}

func installYarn(r *genny.Runner) error {
	// if there's no yarn, install it!
	if _, err := exec.LookPath("yarnpkg"); err == nil {
		return nil
	}
	yargs := []string{"install", "-g", "yarn"}
	return r.Exec(exec.Command("npm", yargs...))
}
