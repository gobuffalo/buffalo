package build

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/packr/v2/jam"
	"github.com/pkg/errors"
)

// Cleanup all of the generated files
func Cleanup(opts *Options) genny.RunFn {
	return func(r *genny.Runner) error {
		defer os.RemoveAll(filepath.Join(opts.Root, "a"))
		if err := jam.Clean(); err != nil {
			return errors.WithStack(err)
		}

		var err error
		opts.rollback.Range(func(k, v interface{}) bool {
			f := genny.NewFileS(k.(string), v.(string))
			r.Logger.Debugf("Rollback: %s", f.Name())
			if err = r.File(f); err != nil {
				return false
			}
			r.Disk.Remove(f.Name())
			return true
		})
		if err != nil {
			return errors.WithStack(err)
		}
		for _, f := range r.Disk.Files() {
			if err := r.Disk.Delete(f.Name()); err != nil {
				return errors.WithStack(err)
			}
		}
		if envy.Mods() {
			if err := r.Exec(exec.Command(genny.GoBin(), "mod", "tidy")); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}
}
