package actions

import (
	"fmt"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/packr/v2"
	"github.com/pkg/errors"
)

type data map[string]interface{}

type presenter struct {
	Name    name.Ident
	Actions []name.Ident
	Helpers data
	Data    data
}

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}
	box := packr.New("github.com/gobuffalo/buffalo/genny/actions/templates", "../actions/templates")

	pres := presenter{
		Name:    name.New(opts.Name),
		Data:    data{},
		Helpers: data{},
	}
	for _, a := range opts.Actions {
		pres.Actions = append(pres.Actions, name.New(a))
	}

	h, err := box.FindString("actions_header.go.tmpl")
	if err != nil {
		return g, errors.WithStack(err)
	}
	a, err := box.FindString("actions.go.tmpl")
	if err != nil {
		return g, errors.WithStack(err)
	}

	f := genny.NewFileS("actions/"+pres.Name.File().String()+".go.tmpl", h+a)
	g.File(f)

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return errors.WithStack(err)
		}
		for _, a := range pres.Actions {
			e := fmt.Sprintf("app.GET(\"/%s/%s\", %s%s)", pres.Name.Underscore(), a.Underscore(), pres.Name.Pascalize(), a.Pascalize())
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

	g.RunFn(func(r *genny.Runner) error {
		f, err := box.FindString("view.html.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}
		for _, a := range pres.Actions {
			pres.Data["action"] = a
			fn := fmt.Sprintf("templates/%s/%s.html.tmpl", pres.Name.Folder(), a.File())
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

	f = genny.NewFileS("actions/"+pres.Name.File().String()+"_test.go.tmpl", h+a)
	g.File(f)

	pres.Data["actions"] = pres.Actions
	pres.Data["name"] = pres.Name
	tmpl := gotools.TemplateTransformer(pres.Data, pres.Helpers)
	g.Transformer(tmpl)
	return g, nil
}
