// +build integration_test

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/dep"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/require"
)

func Test_NewCmd_NoName(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"-f",
	})
	err := c.Execute()
	r.EqualError(err, "you must enter a name for your new application")
}

func Test_NewCmd_InvalidDBType(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"-f",
		"coke",
		"--db-type",
		"a",
	})
	err := c.Execute()
	r.EqualError(err, fmt.Sprintf("unknown dialect a expecting one of %s", strings.Join(pop.AvailableDialects, ", ")))
}

func Test_NewCmd_ForbiddenAppName(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"-f",
		"buffalo",
	})
	err := c.Execute()
	r.EqualError(err, "name buffalo is not allowed, try a different application name")
}

func Test_NewCmd_Nominal(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	gp, err := envy.MustGet("GOPATH")
	r.NoError(err)
	cpath := filepath.Join(gp, "src", "github.com", "gobuffalo")
	tdir, err := ioutil.TempDir(cpath, "testapp")
	r.NoError(err)
	defer os.RemoveAll(tdir)

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	c.SetArgs([]string{
		"new",
		"-f",
		"hello_world",
		"--skip-pop",
		"--skip-webpack",
		"--vcs=none",
	})
	err = c.Execute()
	r.NoError(err)

	r.DirExists(filepath.Join(tdir, "hello_world"))
}

func Test_NewCmd_API(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	gp, err := envy.MustGet("GOPATH")
	r.NoError(err)
	cpath := filepath.Join(gp, "src", "github.com", "gobuffalo")
	tdir, err := ioutil.TempDir(cpath, "testapp")
	r.NoError(err)
	defer os.RemoveAll(tdir)

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	c.SetArgs([]string{
		"new",
		"-f",
		"hello_world",
		"--skip-pop",
		"--api",
		"--vcs=none",
	})
	err = c.Execute()
	r.NoError(err)

	r.DirExists(filepath.Join(tdir, "hello_world"))
}

func Test_NewCmd_WithDep(t *testing.T) {
	envy.Set(envy.GO111MODULE, "off")
	c := RootCmd

	r := require.New(t)
	gp, err := envy.MustGet("GOPATH")
	r.NoError(err)

	newApp := func(rr *require.Assertions) {
		cpath := filepath.Join(gp, "src", "github.com", "gobuffalo")
		tdir, err := ioutil.TempDir(cpath, "testapp")
		rr.NoError(err)
		defer os.RemoveAll(tdir)

		pwd, err := os.Getwd()
		rr.NoError(err)
		os.Chdir(tdir)
		defer os.Chdir(pwd)

		c.SetArgs([]string{
			"new",
			"-f",
			"hello_world",
			"--skip-pop",
			"--skip-webpack",
			"--with-dep",
			"--vcs=none",
			"-v",
		})
		err = c.Execute()
		rr.NoError(err)

		rr.DirExists(filepath.Join(tdir, "hello_world"))
		rr.FileExists(filepath.Join(tdir, "hello_world", "Gopkg.toml"))
		rr.FileExists(filepath.Join(tdir, "hello_world", "Gopkg.lock"))
		rr.DirExists(filepath.Join(tdir, "hello_world", "vendor"))
	}

	// make sure dep installed
	run := genny.WetRunner(context.Background())
	run.WithRun(dep.InstallDep())
	r.NoError(run.Run())

	newApp(r)
}

func Test_NewCmd_WithPopSQLite3(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	gp, err := envy.MustGet("GOPATH")
	r.NoError(err)
	cpath := filepath.Join(gp, "src", "github.com", "gobuffalo")
	tdir, err := ioutil.TempDir(cpath, "testapp")
	r.NoError(err)
	r.NoError(os.MkdirAll(tdir, 0755))
	defer os.RemoveAll(tdir)

	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(tdir)
	defer os.Chdir(pwd)

	c.SetArgs([]string{
		"new",
		"-f",
		"hello_world",
		"--db-type=sqlite3",
		"--skip-webpack",
		"--vcs=none",
		"-v",
	})
	err = c.Execute()
	r.NoError(err)

	r.DirExists(filepath.Join(tdir, "hello_world"))
	r.FileExists(filepath.Join(tdir, "hello_world", "database.yml"))
}
