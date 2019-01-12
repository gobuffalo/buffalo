package actions

import (
	"fmt"
	"strings"

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
	Options *Options
}

var box = packr.New("github.com/gobuffalo/buffalo/genny/actions/templates", "../actions/templates")

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	g.RunFn(construct(opts))
	return g, nil
}

func buildTemplates(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		f, err := box.FindString("view.html.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}
		for _, a := range pres.Actions {
			pres.Data["action"] = a
			fn := fmt.Sprintf("templates/%s/%s.html.tmpl", pres.Name.Folder(), a.File())
			xf := genny.NewFileS(fn, f)
			xf, err = transform(pres, xf)
			if err != nil {
				return errors.WithStack(err)
			}
			if err := r.File(xf); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}
}

func updateApp(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return errors.WithStack(err)
		}
		var lines []string
		body := f.String()
		for _, a := range pres.Actions {
			e := fmt.Sprintf("app.GET(\"/%s/%s\", %s%s)", pres.Name.Underscore(), a.Underscore(), pres.Name.Pascalize(), a.Pascalize())
			if !strings.Contains(body, e) {
				lines = append(lines, e)
			}
		}
		f, err = gotools.AddInsideBlock(f, "app == nil", strings.Join(lines, "\n\t\t"))
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	}
}

func buildActions(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		fn := fmt.Sprintf("actions/%s.go", pres.Name.File())
		xf, err := r.FindFile(fn)
		if err != nil {
			return buildNewActions(fn, pres)(r)
		}
		if err := appendActions(xf, pres)(r); err != nil {
			return err
		}

		return nil
	}
}

func appendActions(f genny.File, pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		body := f.String()
		for _, ac := range pres.Options.Actions {
			a := name.New(ac)
			x := fmt.Sprintf("func %s%s", pres.Name.Pascalize(), a.Pascalize())
			if strings.Contains(body, x) {
				continue
			}
			pres.Actions = append(pres.Actions, a)
		}

		a, err := box.FindString("actions.go.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}

		f = genny.NewFileS(f.Name()+".tmpl", f.String()+a)

		f, err = transform(pres, f)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	}
}

func buildNewActions(fn string, pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		for _, a := range pres.Options.Actions {
			pres.Actions = append(pres.Actions, name.New(a))
		}

		h, err := box.FindString("actions_header.go.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}
		a, err := box.FindString("actions.go.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}

		f := genny.NewFileS(fn+".tmpl", h+a)

		f, err = transform(pres, f)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	}
}

func construct(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		pres := &presenter{
			Name:    name.New(opts.Name),
			Data:    data{},
			Helpers: data{},
			Options: opts,
		}

		if err := buildActions(pres)(r); err != nil {
			return err
		}

		if err := buildTests(pres)(r); err != nil {
			return err
		}

		if err := updateApp(pres)(r); err != nil {
			return errors.WithStack(err)
		}

		if err := buildTemplates(pres)(r); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}
}

func transform(pres *presenter, f genny.File) (genny.File, error) {
	pres.Data["actions"] = pres.Actions
	pres.Data["name"] = pres.Name
	t := gotools.TemplateTransformer(pres.Data, pres.Helpers)
	return t.Transform(f)
}

func buildTests(pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		fn := fmt.Sprintf("actions/%s_test.go", pres.Name.File())
		xf, err := r.FindFile(fn)
		if err != nil {
			return buildNewTests(fn, pres)(r)
		}
		return appendTests(xf, pres)(r)
	}
}

func appendTests(f genny.File, pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		a, err := box.FindString("test.go.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}
		f := genny.NewFileS(f.Name()+".tmpl", f.String()+a)
		f, err = transform(pres, f)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	}
}

func buildNewTests(fn string, pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		h, err := box.FindString("tests_header.go.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}
		a, err := box.FindString("test.go.tmpl")
		if err != nil {
			return errors.WithStack(err)
		}

		f := genny.NewFileS(fn+".tmpl", h+a)

		f, err = transform(pres, f)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	}
}
