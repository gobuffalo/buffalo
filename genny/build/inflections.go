package build

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func apkg(opts *Options) genny.RunFn {
	rns := []genny.RunFn{
		copyInflections,
		copyDatabase,
	}

	return func(r *genny.Runner) error {
		for _, rn := range rns {
			if err := rn(r); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}
}

func copyDatabase(r *genny.Runner) error {
	f, err := r.FindFile("database.yml")
	if err != nil {
		// it's ok to not have this file
		return nil
	}
	bb := &bytes.Buffer{}
	io.Copy(bb, f)

	dgo := &bytes.Buffer{}
	dgo.WriteString("package a\n")
	dgo.WriteString(fmt.Sprintf("var DB_CONFIG = `%s`", bb.String()))
	return r.File(genny.NewFile("a/database.go", dgo))
}

func copyInflections(r *genny.Runner) error {
	f, err := r.FindFile("inflections.json")
	if err != nil {
		// it's ok to not have this file
		return nil
	}
	return r.File(genny.NewFile("a/inflections.json", f))
}
