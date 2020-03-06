package build

import (
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

func buildDeps(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	if envy.Mods() {
		return g, nil
	}

	// mount the go get runner
	tf := opts.App.BuildTags(opts.Environment, opts.Tags...)
	if len(tf) > 0 {
		tf = append([]string{"-tags"}, tf.String())
	}
	g.Command(gogen.Get("./...", tf...))
	return g, nil
}
