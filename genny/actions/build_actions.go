package actions

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
)

// buildActions is the top level action builder
// it determines whether to build a new actions/foo.go file
// or append to an existing one
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

// buildNewActions builds a brand new actions/foo.go file
// and files it with actions
func buildNewActions(fn string, pres *presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		for _, a := range pres.Options.Actions {
			pres.Actions = append(pres.Actions, name.New(a))
		}

		h, err := box.FindString("actions_header.go.tmpl")
		if err != nil {
			return err
		}
		a, err := box.FindString("actions.go.tmpl")
		if err != nil {
			return err
		}

		f := genny.NewFileS(fn+".tmpl", h+a)

		f, err = transform(pres, f)
		if err != nil {
			return err
		}
		return r.File(f)
	}
}

// appendActions appends *only* actions that don't exist in
// actions/foo.go. if the action already exists it is not touched.
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
			return err
		}

		f = genny.NewFileS(f.Name()+".tmpl", f.String()+a)

		f, err = transform(pres, f)
		if err != nil {
			return err
		}
		return r.File(f)
	}
}
