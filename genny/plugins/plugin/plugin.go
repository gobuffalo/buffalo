package plugin

import (
	"path/filepath"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/licenser/genny/licenser"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
	"github.com/gobuffalo/release/genny/initgen"
	"github.com/gobuffalo/release/genny/makefile"
)

// New generates a plugin project
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, err
	}

	g := genny.New()
	g.Box(packr.New("buffalo:genny:plugin", "../plugin/templates"))
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("-short-", opts.ShortName))
	g.Transformer(genny.Dot())
	gg.Add(g)

	lopts := &licenser.Options{
		Author: opts.Author,
		Name:   opts.License,
	}

	g, err := licenser.New(lopts)
	if err != nil {
		return gg, err
	}
	gg.Add(g)

	ig, err := initgen.New(&initgen.Options{
		Version:     "v0.0.1",
		VersionFile: filepath.Join(opts.ShortName, "version.go"),
		Options: &makefile.Options{
			MainFile: "main.go",
			Root:     opts.Root,
		},
	})
	if err != nil {
		return gg, err
	}
	gg.Merge(ig)

	return gg, nil
}
