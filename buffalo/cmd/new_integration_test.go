// +build integration_test

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/dep"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/require"
)

func Test_NewCmd_NoName(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"-v",
	})
	err := c.Execute()
	r.EqualError(err, "you must enter a name for your new application")
}

func Test_NewCmd_InvalidDBType(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"-v",
		"coke",
		"--db-type",
		"x",
	})
	err := c.Execute()
	r.EqualError(err, fmt.Sprintf("unknown dialect \"x\" expecting one of %s", strings.Join(pop.AvailableDialects, ", ")))
}

func Test_NewCmd_ForbiddenAppName(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"-v",
		"buffalo",
	})
	err := c.Execute()
	r.EqualError(err, "name buffalo is not allowed, try a different application name")
}

func Test_NewCmd_Nominal(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	err := withDir(func(dir string) {
		c.SetArgs([]string{
			"new",
			"--skip-pop",
			"--skip-webpack",
			"--vcs=none",
			"-v",
			"hello_world",
		})
		err := c.Execute()
		r.NoError(err)
		r.DirExists(filepath.Join(dir, "hello_world"))
	})
	r.NoError(err)

}

func Test_NewCmd_API(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	err := withDir(func(dir string) {
		c.SetArgs([]string{
			"new",
			"--skip-pop",
			"--api",
			"--vcs=none",
			"-v",
			"hello_world",
		})
		err := c.Execute()
		r.NoError(err)

		r.DirExists(filepath.Join(dir, "hello_world"))
	})

	r.NoError(err)
}

func Test_NewCmd_WithDep(t *testing.T) {
	r := require.New(t)
	err := gomods.Disable(func() error {
		c := RootCmd
		newApp := func(rr *require.Assertions) error {
			return withDir(func(dir string) {
				c.SetArgs([]string{
					"new",
					"--skip-pop",
					"--skip-webpack",
					"--with-dep",
					"--vcs=none",
					"-v",
					"hello_world",
				})
				err := c.Execute()
				rr.NoError(err)

				rr.DirExists(filepath.Join(dir, "hello_world"))
				rr.FileExists(filepath.Join(dir, "hello_world", "Gopkg.toml"))
				rr.FileExists(filepath.Join(dir, "hello_world", "Gopkg.lock"))
				rr.DirExists(filepath.Join(dir, "hello_world", "vendor"))
			})
		}

		// make sure dep installed
		run := genny.WetRunner(context.Background())
		run.WithRun(dep.InstallDep())
		r.NoError(run.Run())

		return newApp(r)
	})
	r.NoError(err)
}

func Test_NewCmd_WithPopSQLite3(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	err := withDir(func(dir string) {

		c.SetArgs([]string{
			"new",
			"--db-type=sqlite3",
			"--skip-webpack",
			"--vcs=none",
			"-v",
			"hello_world",
		})
		err := c.Execute()
		r.NoError(err)

		r.DirExists(filepath.Join(dir, "hello_world"))
		r.FileExists(filepath.Join(dir, "hello_world", "database.yml"))
	})
	r.NoError(err)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func withDir(fn func(string)) error {
	gp, err := envy.MustGet("GOPATH")
	if err != nil {
		return err
	}
	cpath := filepath.Join(gp, "src", "github.com", "gobuffalo")
	tdir, err := ioutil.TempDir(cpath, fmt.Sprint(rand.Int()))
	if err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	fn(tdir)
	os.RemoveAll(tdir)
	return nil
}
