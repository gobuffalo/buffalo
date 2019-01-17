package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/buffalo/buffalo/cmd"
	"github.com/gobuffalo/envy"
	"github.com/markbates/safe"
)

func call(args []string, fn func(dir string)) error {
	gp, err := envy.MustGet("GOPATH")
	if err != nil {
		return err
	}
	if fn == nil {
		if err := exec(args); err != nil {
			return err
		}
	}
	cpath := filepath.Join(gp, "src", "github.com", "gobuffalo")
	tdir, err := ioutil.TempDir(cpath, "testapp")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tdir)

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(tdir)
	defer os.Chdir(pwd)

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
