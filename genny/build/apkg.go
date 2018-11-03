package build

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

func apkg(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	g.RunFn(copyInflections)
	g.RunFn(copyDatabase)

	return g, nil
}

func copyDatabase(r *genny.Runner) error {
	defer func() {
		r.Disk.Remove("database.yml")
	}()

	bb := &bytes.Buffer{}

	f, err := r.FindFile("database.yml")
	if err == nil {
		io.Copy(bb, f)
	}

	dgo := &bytes.Buffer{}
	dgo.WriteString("package a\n")
	dgo.WriteString(fmt.Sprintf("var DB_CONFIG = `%s`", bb.String()))
	return r.File(genny.NewFile("a/database.go", dgo))
}

func copyInflections(r *genny.Runner) error {
	defer func() {
		r.Disk.Remove("inflections.json")
	}()
	f, err := r.FindFile("inflections.json")
	if err != nil {
		// it's ok to not have this file
		return nil
	}
	return r.File(genny.NewFile("a/inflections.json", f))
}
