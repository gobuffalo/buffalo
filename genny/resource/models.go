package resource

import (
	"os/exec"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
)

func modelCommand(model name.Ident, opts *Options) *exec.Cmd {
	args := opts.Attrs.Slice()
	mn := model.Singularize().Underscore().String()
	args = append([]string{"pop", "g", "model", mn}, args...)

	if opts.SkipMigration {
		args = append(args, "--skip-migration")
	}

	return exec.Command("buffalo-pop", args...)
}

func installPop(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		if opts.SkipModel {
			return nil
		}
		if _, err := r.LookPath("buffalo-pop"); err != nil {
			c := gogen.Get("github.com/gobuffalo/buffalo-pop")
			if err := r.Exec(c); err != nil {
				return err
			}
		}
		return r.Exec(modelCommand(name.New(opts.Model), opts))
	}
}
