package integration

import (
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/buffalo/cmd"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/packr/v2/jam"
	"github.com/markbates/safe"
)

func call(args []string, fn func(dir string)) error {
	jam.Clean()
	defer jam.Clean()
	ogp, err := envy.MustGet("GOPATH")
	defer envy.MustSet("GOPATH", ogp)
	gp := os.TempDir()
	err = envy.MustSet("GOPATH", gp)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	if fn == nil {
		if err := exec(args); err != nil {
			return err
		}
	}
	tdir := filepath.Join(gp, "src", "github.com", "gobuffalo", "testapp")
	defer os.RemoveAll(tdir)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(tdir); err != nil {
		return err
	}
	if err := os.MkdirAll(tdir, 0755); err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(tdir)
	os.Setenv("PWD", tdir)
	defer os.Chdir(pwd)
	defer os.Setenv("PWD", pwd)

	if err := exec(args); err != nil {
		return err
	}
	return safe.Run(func() {
		fn(tdir)
	})
}

func exec(args []string) error {
	c := cmd.RootCmd
	c.SetArgs(args)
	return c.Execute()
}
