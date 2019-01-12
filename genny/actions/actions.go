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

var box = packr.New("github.com/gobuffalo/buffalo/genny/actions/templates", "../actions/templates")

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	g.RunFn(construct(opts))
	return g, nil
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
