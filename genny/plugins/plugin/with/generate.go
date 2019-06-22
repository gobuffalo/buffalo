package with

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo/genny/plugins/plugin"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/genny/new"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush"
	"github.com/gobuffalo/plushgen"
)

// GenerateCmd generates a plugin project with go mods
func GenerateCmd(opts *plugin.Options) (*genny.Group, error) {
	gg := &genny.Group{}
	if err := opts.Validate(); err != nil {
		return gg, err
	}

	g := genny.New()
	box := packr.New("./generate/templates", "./generate/templates")
	if err := g.Box(box); err != nil {
		return gg, err
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	g.Transformer(plushgen.Transformer(ctx))

	g.Transformer(genny.Replace("-shortName-", opts.ShortName))
	g.Transformer(genny.Dot())

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("cmd/available.go")
		if err != nil {
			return err
		}
		const g = `Available.Add("generate", generateCmd)`
		const m = `Available.Mount(rootCmd)`
		body := strings.Replace(f.String(), m, fmt.Sprintf("\t%s\n%s", g, m), 1)
		return r.File(genny.NewFile(f.Name(), strings.NewReader(body)))
	})

	gg.Add(g)

	g, err := new.New(&new.Options{
		Name:   opts.ShortName,
		Prefix: "genny",
	})
	if err != nil {
		return gg, err
	}
	gg.Add(g)

	return gg, nil
}
