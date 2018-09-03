package resource

import (
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

// New resource generator for pop
func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	if err := opts.Validate(); err != nil {
		return gg, errors.WithStack(err)
	}
	return gg, nil
}

func modelCommand(opts *Options) *exec.Cmd {
	args := []string{"pop", "g", "model", opts.Attrs.Name.String()}
	for _, a := range opts.Attrs.Attrs {
		args = append(args, a.String())
	}

	if opts.SkipMigration {
		args = append(args, "--skip-migration")
	}

	return exec.Command("buffalo", args...)
}
