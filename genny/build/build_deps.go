package build

import (
	"os/exec"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/v2"
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
	args := []string{"get"}
	args = append(args, tf...)
	args = append(args, "./...")
	g.Command(exec.Command("go", args...))
	return g, nil
}
