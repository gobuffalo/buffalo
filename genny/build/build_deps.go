package build

import (
	"github.com/gobuffalo/depgen"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/pkg/errors"
)

func buildDeps(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if opts.App.WithDep {
		// mount the dep generator
		dg, err := depgen.Ensure(false)
		if err != nil {
			return g, errors.WithStack(err)
		}
		g.Merge(dg)
	} else {
		// mount the go get runner
		tf := opts.App.BuildTags(opts.Environment, opts.Tags...)
		if len(tf) > 0 {
			tf = append([]string{"-tags"}, tf.String())
		}
		g.Command(gogen.Get("./...", tf...))
	}
	return g, nil
}
