package resource

import (
	"text/template"

	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
)

// New resource generator
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if !opts.SkipTemplates {
		core := packr.New("github.com/gobuffalo/buffalo/@v0.15.4/genny/resource/templates/core", "../resource/templates/core")

		if err := g.Box(core); err != nil {
			return g, err
		}
	}

	var abox packd.Box
	if opts.SkipModel {
		abox = packr.New("github.com/gobuffalo/buffalo/@v0.15.4/genny/resource/templates/standard", "../resource/templates/standard")
	} else {
		abox = packr.New("github.com/gobuffalo/buffalo/@v0.15.4/genny/resource/templates/use_model", "../resource/templates/use_model")
	}

	if err := g.Box(abox); err != nil {
		return g, err
	}

	pres := presenter{
		App:   opts.App,
		Name:  name.New(opts.Name),
		Model: name.New(opts.Model),
		Attrs: opts.Attrs,
	}
	x := pres.Name.Resource().File().String()
	folder := pres.Name.Folder().Pluralize().String()
	g.Transformer(genny.Replace("resource-name", x))
	g.Transformer(genny.Replace("resource-use_model", x))
	g.Transformer(genny.Replace("folder-name", folder))

	data := map[string]interface{}{
		"opts":    pres,
		"actions": actions(opts),
		"folder":  folder,
	}
	helpers := template.FuncMap{
		"camelize": func(s string) string {
			return flect.Camelize(s)
		},
	}
	g.Transformer(gogen.TemplateTransformer(data, helpers))

	g.RunFn(installPop(opts))

	g.RunFn(addResource(pres))
	return g, nil
}
