package standard

import (
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr/v2"
)

// New generator for creating basic asset files
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	g.Box(packr.New("buffalo:genny:assets:standard", "../standard/templates"))

	data := map[string]interface{}{}
	h := template.FuncMap{}
	t := gotools.TemplateTransformer(data, h)
	g.Transformer(t)

	return g, nil
}
