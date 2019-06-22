package core

import (
	"github.com/gobuffalo/buffalo/genny/docker"
	pop "github.com/gobuffalo/buffalo-pop/genny/newapp"
	"github.com/gobuffalo/buffalo/genny/ci"
	"github.com/gobuffalo/buffalo/genny/plugins/install"
	"github.com/gobuffalo/buffalo/genny/refresh"
	"github.com/gobuffalo/buffalo/plugins/plugdeps"
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/depgen"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/gogen/gomods"
	"github.com/gobuffalo/meta"
	"github.com/pkg/errors"
)

// New generator for creating a Buffalo application
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	// add the root generator
	g, err := rootGenerator(opts)
	if err != nil {
		return gg, err
	}
	gg.Add(g)

	app := opts.App

	if app.WithModules {
		g, err := gomods.Init(app.PackagePkg, app.Root)
		if err != nil {
			return gg, err
		}
		g.Command(gogen.Get("github.com/gobuffalo/buffalo@" + runtime.Version))
		g.Command(gogen.Get("./..."))

		gg.Add(g)
	}

	plugs, err := plugdeps.List(app)
	if err != nil && (errors.Cause(err) != plugdeps.ErrMissingConfig) {
		return nil, err
	}

	if opts.Docker != nil {
		// add the docker generator
		g, err = docker.New(opts.Docker)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	if opts.Pop != nil {
		// add the pop generator
		gg2, err := pop.New(opts.Pop)
		if err != nil {
			return gg, err
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
			return gg, err
		}
		gg.Add(g)
	}

	if opts.Refresh != nil {
		g, err = refresh.New(opts.Refresh)
		if err != nil {
			return gg, err
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
		return gg, err
	}
	gg.Merge(ig)

	// DEP/MODS/go get should be last
	if app.WithDep {
		// init dep
		di, err := depgen.Init("", false)
		if err != nil {
			return gg, err
		}
		gg.Add(di)
	}

	if !app.WithDep && !app.WithModules {
		g := genny.New()
		g.Command(gogen.Get("./...", "-t"))
		gg.Add(g)
	}

	if app.WithModules {
		g, err := gomods.Tidy(app.Root, false)
		if err != nil {
			return gg, err
		}
		gg.Add(g)
	}

	return gg, nil
}
