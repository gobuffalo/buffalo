package resource

import (
	"fmt"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
)

func addResource(pres presenter) genny.RunFn {
	return func(r *genny.Runner) error {
		f, err := r.FindFile("actions/app.go")
		if err != nil {
			return errors.WithStack(err)
		}
		stmt := fmt.Sprintf("app.Resource(\"/%s\", %sResource{})", pres.Name.URL(), pres.Name.Resource())
		f, err = gotools.AddInsideBlock(f, "if app == nil {", stmt)
		if err != nil {
			return errors.WithStack(err)
		}
		return r.File(f)
	}
}

func actions(opts *Options) []name.Ident {
	actions := []name.Ident{
		name.New("list"),
		name.New("show"),
		name.New("create"),
		name.New("update"),
		name.New("destroy"),
	}

	if opts.App.AsWeb {
		actions = append(actions, name.New("new"), name.New("edit"))
	}

	return actions
}
