package refresh

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

// New generator to generate refresh templates
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}
	g.Box(packr.NewBox("../refresh/templates"))

	ctx := plush.NewContext()
	ctx.Set("app", opts.App)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Dot())
	return g, nil
}
