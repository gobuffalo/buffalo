package build

import (
	"os/exec"

	"github.com/gobuffalo/genny/v2"
)

func apkg(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, err
	}

	g.RunFn(copyInflections)
	g.RunFn(copyDatabase)
	g.RunFn(addDependencies)

	return g, nil
}

func addDependencies(r *genny.Runner) error {
	return r.Exec(exec.Command("go", "get", "-d"))
}

func copyDatabase(r *genny.Runner) error {
	defer func() {
		r.Disk.Remove("database.yml")
	}()
	f, err := r.FindFile("database.yml")
	if err != nil {
		f, err = r.FindFile("config/database.yml")
		if err != nil {
			// it's ok to not have this file
			return nil
		}
	}
	return r.File(genny.NewFile("a/database.yml", f))
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
