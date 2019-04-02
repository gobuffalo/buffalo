package core

import (
	"os"
	"strings"

	"github.com/gobuffalo/genny"
)

func validateInGoPath(srcDirs []string) genny.RunFn {
	return func(r *genny.Runner) error {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		var found bool
		for _, src := range srcDirs {
			if strings.Contains(pwd, src) {
				found = true
				break
			}
		}
		if !found {
			return ErrNotInGoPath
		}
		return nil
	}
}
