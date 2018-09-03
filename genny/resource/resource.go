package resource

import (
	"fmt"
	"path/filepath"

	"github.com/gobuffalo/buffalo/genny/resource/pop/resource"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
)

// New resource generator group
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}
	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}

	g := genny.New()
	g.RunFn(func(r *genny.Runner) error {
		return mapResource(r, opts)
	})
	gg.Add(g)

	pop, err := resource.New(&resource.Options{
		App:           opts.App,
		Attrs:         opts.Attrs,
		SkipMigration: opts.SkipMigration,
		SkipModel:     opts.SkipModel,
		SkipTemplates: opts.SkipTemplates,
		UseModel:      opts.UseModel,
	})
	if err != nil {
		return gg, errors.WithStack(err)
	}

	gg.Merge(pop)

	return gg, nil
}

func mapResource(r *genny.Runner, opts *Options) error {
	f, err := r.FindFile(filepath.Join("actions", "app.go"))
	if err != nil {
		return errors.WithStack(err)
	}
	n := opts.Attrs.Name
	f, err = gotools.AddInsideBlock(f, "if app == nil {", fmt.Sprintf("app.Resource(\"/%s\", %sResource{})", n.URL(), n.Resource()))
	if err != nil {
		return errors.WithStack(err)
	}
	return r.File(f)
}
