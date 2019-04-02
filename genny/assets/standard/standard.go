package standard

import (
	"strings"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

// New generator for creating basic asset files
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	g.Box(packr.New("buffalo:genny:assets:standard", "../standard/templates"))

	data := map[string]interface{}{}
	h := template.FuncMap{}
	t := gogen.TemplateTransformer(data, h)
	g.Transformer(t)

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("templates/application.html")
		if err != nil {
			return errors.WithStack(err)
		}

		s := strings.Replace(f.String(), "</title>", "</title>\n"+bs4, 1)
		return r.File(genny.NewFileS(f.Name(), s))
	})

	return g, nil
}

const bs4 = `<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-9gVQ4dYFwwWSjIDZnLEWnxCjeSWFphJiwGPXr1jddIhOegiu1FwO5qRGvFXOdJZ4" crossorigin="anonymous">`
