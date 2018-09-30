package build

import (
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/dep"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
)

func (b *Builder) installBuildDeps() error {
	run := genny.WetRunner(b.ctx)
	if b.WithDep {
		if err := run.WithNew(dep.Ensure(b.Debug)); err != nil {
			return errors.WithStack(err)
		}
	} else {
		run.WithRun(gotools.Get("./..."))
	}
	return run.Run()
}
