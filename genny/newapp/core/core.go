package core

import (
	"os/exec"

	"github.com/gobuffalo/buffalo-docker/genny/docker"
	"github.com/gobuffalo/buffalo-plugins/genny/install"
	"github.com/gobuffalo/buffalo-plugins/plugins/plugdeps"
	pop "github.com/gobuffalo/buffalo-pop/genny/newapp"
	"github.com/gobuffalo/buffalo/genny/ci"
	"github.com/gobuffalo/buffalo/genny/refresh"
	"github.com/gobuffalo/buffalo/genny/vcs"
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/dep"
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
)

// New generator for creating a Buffalo application
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	// add the root generator
	g, err := rootGenerator(opts)
	if err != nil {
		return gg, errors.WithStack(err)
	}
	gg.Add(g)

	app := opts.App

	plugs, err := plugdeps.List(app)
	if err != nil && (errors.Cause(err) != plugdeps.ErrMissingConfig) {
		return nil, errors.WithStack(err)
	}

	if opts.Docker != nil {
		// add the docker generator
		g, err = docker.New(opts.Docker)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(g)
	}

	if opts.Pop != nil {
		// add the pop generator
		gg2, err := pop.New(opts.Pop)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Merge(gg2)

		// add the plugin
		plugs.Add(plugdeps.Plugin{
			Binary: "buffalo-pop",
			GoGet:  "github.com/gobuffalo/buffalo-pop",
		})
	}

	if opts.CI != nil {
		// add the CI generator
		g, err = ci.New(opts.CI)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(g)
	}

	if opts.Refresh != nil {
		g, err = refresh.New(opts.Refresh)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(g)
	}

	// ---

	// install all of the plugins
	iopts := &install.Options{
		App:     app,
		Plugins: plugs.List(),
	}
	if app.WithSQLite {
		iopts.Tags = meta.BuildTags{"sqlite"}
	}

	ig, err := install.New(iopts)
	if err != nil {
		return gg, errors.WithStack(err)
	}
	gg.Merge(ig)

	// DEP/MODS/go get should be last
	if app.WithDep {
		// init dep
		di, err := dep.Init("", false)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(di)
	}

	if app.WithModules {
		g := genny.New()
		g.Command(exec.Command(genny.GoBin(), "get", "github.com/gobuffalo/buffalo@"+runtime.Version))
		g.Command(exec.Command(genny.GoBin(), "get"))
		g.Command(exec.Command(genny.GoBin(), "mod", "tidy"))
		gg.Add(g)
	}

	if !app.WithDep && !app.WithModules {
		g := genny.New()
		g.Command(exec.Command(genny.GoBin(), "get", "-t", "./..."))
		gg.Add(g)
	}

	// setup VCS last
	if opts.VCS != nil {
		// add the VCS generator
		g, err = vcs.New(opts.VCS)
		if err != nil {
			return gg, errors.WithStack(err)
		}
		gg.Add(g)
	}

	return gg, nil
}
