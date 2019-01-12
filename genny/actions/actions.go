package actions

import (
	"fmt"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}
	box := packr.New("github.com/gobuffalo/buffalo/genny/actions/templates", "../actions/templates")

	actions := []name.Ident{}
	for _, a := range opts.Actions {
		actions = append(actions, name.New(a))
	}

	n := name.New(opts.Name)

	h, err := box.FindString("actions_header.go.tmpl")
	if err != nil {
		return g, errors.WithStack(err)
	}
	a, err := box.FindString("actions.go.tmpl")
	if err != nil {
		return g, errors.WithStack(err)
	}

	f := genny.NewFileS("actions/"+n.File().String()+".go.tmpl", h+a)
	g.File(f)

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return errors.WithStack(err)
		}
		for _, a := range actions {
			e := fmt.Sprintf("app.GET(\"/%s/%s\", %s%s)", n.Underscore(), a.Underscore(), n.Pascalize(), a.Pascalize())
			f, err := gotools.AddInsideBlock(f, "app == nil", e)
			if err != nil {
				return errors.WithStack(err)
			}
			if err := r.File(f); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	help := map[string]interface{}{}
	data := map[string]interface{}{
		"name":    n,
		"actions": actions,
	}

	g.RunFn(func(r *genny.Runner) error {
		f, err := box.FindString("view.html.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}
		for _, a := range actions {
			data["action"] = a
			fn := fmt.Sprintf("templates/%s/%s.html.tmpl", n.Folder(), a.File())
			xf := genny.NewFileS(fn, f)
			if err := r.File(xf); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	h, err = box.FindString("tests_header.go.tmpl")
	if err != nil {
		return g, errors.WithStack(err)
	}
	a, err = box.FindString("test.go.tmpl")
	if err != nil {
		return g, errors.WithStack(err)
	}

	f = genny.NewFileS("actions/"+n.File().String()+"_test.go.tmpl", h+a)
	g.File(f)

	tmpl := gotools.TemplateTransformer(data, help)
	g.Transformer(tmpl)
	return g, nil
}
