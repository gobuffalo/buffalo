package docker

import (
	"strings"
	"text/template"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/runtime"
)

var boxes = map[string]packr.Box{
	"standard": packr.NewBox("./standard/templates"),
	"multi":    packr.NewBox("./multi/templates"),
}

// New generator for Dockerfiles
func New(opts *Options) (*genny.Generator, error) {
	if err := validateOptions(opts); err != nil {
		return nil, errors.WithStack(err)
	}
	g := genny.New()
	box, ok := boxes[opts.Style]
	if !ok {
		return g, errors.Errorf("unknown Docker style: %s", opts.Style)
	}
	g.Box(box)

	data := map[string]interface{}{
		"opts": opts,
	}
	t := gotools.TemplateTransformer(data, template.FuncMap{})
	g.Transformer(t)
	g.Transformer(genny.Replace("-dot-", "."))
	return g, nil
}

func validateOptions(opts *Options) error {
	if strings.ToLower(opts.Style) == "none" {
		return errors.New("style was none - generator should not be used")
	}
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}
	if opts.Version == "" {
		opts.Version = runtime.Version
	}
	if opts.Style == "" {
		opts.Style = "multi"
	}
	opts.AsWeb = opts.App.WithWebpack
	if _, ok := boxes[opts.Style]; !ok {
		return errors.Errorf("unknown Docker style: %s", opts.Style)
	}

	return nil
}
