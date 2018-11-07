package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func validateInGoPath(srcDirs []string) genny.RunFn {
	fmt.Println("### srcDirs ->", srcDirs)
	return func(r *genny.Runner) error {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		var found bool
		for _, src := range srcDirs {
			if strings.HasPrefix(pwd, src) {
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
