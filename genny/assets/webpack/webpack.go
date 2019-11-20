package webpack

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen"
	"github.com/gobuffalo/packr/v2"
)

// BinPath is the path to the local install of webpack
var BinPath = func() string {
	s := filepath.Join("node_modules", ".bin", "webpack")
	if runtime.GOOS == "windows" {
		s += ".cmd"
	}
	return s
}()

// Templates used for generating webpack
// (exported mostly for the "fix" command)
var Templates = packr.New("github.com/gobuffalo/buffalo/genny/assets/webpack", "../webpack/templates")

// New generator for creating webpack asset files
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.RunFn(func(r *genny.Runner) error {
		if opts.App.WithYarn {
			if _, err := r.LookPath("yarnpkg"); err == nil {
				return nil
			}
			// If yarn is not installed, it still can be installed with npm.
		}
		if _, err := r.LookPath("npm"); err != nil {
			return fmt.Errorf("could not find npm executable")
		}
		return nil
	})

	g.Box(Templates)

	data := map[string]interface{}{
		"opts": opts,
	}
	t := gogen.TemplateTransformer(data, gogen.TemplateHelpers)
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
			return err
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
	if _, err := r.LookPath("yarnpkg"); err == nil {
		return nil
	}
	yargs := []string{"install", "-g", "yarn"}
	return r.Exec(exec.Command("npm", yargs...))
}
