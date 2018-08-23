package grift

import (
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/markbates/inflect"
	"github.com/pkg/errors"
)

// New generator to create a grift task
func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()
	if len(opts.Args) == 0 {
		return g, errors.New("you need to provide a name for the grift task")
	}

	opts.Namespaced = strings.Contains(opts.Args[0], ":")

	for _, n := range strings.Split(opts.Args[0], ":") {
		opts.Parts = append(opts.Parts, inflect.Name(n))
	}
	opts.Name = opts.Parts[len(opts.Parts)-1]

	data := map[string]interface{}{
		"opts": opts,
	}
	t := gotools.TemplateTransformer(data, template.FuncMap{})

	g.RunFn(func(r *genny.Runner) error {
		header := tmplHeader
		path := filepath.Join("grifts", opts.Name.File()+".go.tmpl")
		if f, err := r.FindFile(path); err == nil {
			header = f.String()
		}
		f := genny.NewFile(path, strings.NewReader(header+tmplBody))
		f, err := t.Transform(f)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	})
	return g, nil
}
