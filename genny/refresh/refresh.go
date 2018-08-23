package refresh

import (
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
)

// New generator to generate refresh templates
func New(app meta.App) (*genny.Generator, error) {
	g := genny.New()
	g.Box(packr.NewBox("../refresh/templates"))

	ctx := plush.NewContext()
	ctx.Set("app", app)
	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(genny.Replace("-dot-", "."))
	return g, nil
}
