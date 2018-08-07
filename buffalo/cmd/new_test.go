package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/require"
)

func Test_Bootstrap4_Default(t *testing.T) {
	r := require.New(t)
	f, err := newCmd.Flags().GetInt("bootstrap")
	r.NoError(err)
	r.Equal(4, f)
}

func Test_NewCmd_NoName(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
	})
	err := c.Execute()
	r.EqualError(err, "you must enter a name for your new application")
}

func Test_NewCmd_InvalidDBType(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"coke",
		"--db-type",
		"a",
	})
	err := c.Execute()
	r.EqualError(err, fmt.Sprintf("Unknown db-type a expecting one of %s", strings.Join(pop.AvailableDialects, ", ")))
}

func Test_NewCmd_ForbiddenAppName(t *testing.T) {
	r := require.New(t)
	c := RootCmd
	c.SetArgs([]string{
		"new",
		"buffalo",
	})
	err := c.Execute()
	r.EqualError(err, "name buffalo is not allowed, try a different application name")
}

func Test_NewCmd_Nominal(t *testing.T) {
	r := require.New(t)
	c := RootCmd

	tdir, err := ioutil.TempDir("", "gopath")
	r.NoError(err)
	defer os.RemoveAll(tdir)
	gp, err := envy.MustGet("GOPATH")
	r.NoError(err)
	defer envy.MustSet("GOPATH", gp)
	err = envy.MustSet("GOPATH", tdir)
	r.NoError(err)

	cpath := path.Join(tdir, "src", "github.com", "gobuffalo")
	err = os.MkdirAll(cpath, 0700)
	r.NoError(err)
	pwd, err := os.Getwd()
	r.NoError(err)
	os.Chdir(cpath)
	defer os.Chdir(pwd)

	c.SetArgs([]string{
		"new",
		"hello_world",
		"--skip-pop",
		"--skip-webpack",
		"--vcs=none",
	})
	err = c.Execute()
	r.NoError(err)
}
