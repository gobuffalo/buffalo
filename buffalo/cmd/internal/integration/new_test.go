// +build integration_test

package integration

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/require"
)

func Test_NewCmd_NoName(t *testing.T) {
	err := call([]string{"new"}, nil)
	r := require.New(t)
	r.Error(err)
	r.EqualError(err, "you must enter a name for your new application")
}

func Test_NewCmd_InvalidDBType(t *testing.T) {
	args := []string{
		"new",
		"coke",
		"--db-type",
		"a",
	}
	err := call(args, nil)
	r := require.New(t)
	r.Error(err)
	r.EqualError(err, fmt.Sprintf(`unknown dialect "a" expecting one of %s`, strings.Join(pop.AvailableDialects, ", ")))
}

func Test_NewCmd_ForbiddenAppName(t *testing.T) {
	args := []string{
		"new",
		"buffalo",
	}
	err := call(args, nil)
	r := require.New(t)
	r.Error(err)
	r.EqualError(err, "name buffalo is not allowed, try a different application name")
}

func Test_NewCmd_Nominal(t *testing.T) {
	r := require.New(t)
	args := []string{
		"new",
		"hello_world",
		"--skip-pop",
		"--skip-webpack",
		"--vcs=none",
	}
	err := call(args, func(tdir string) {
		r.DirExists(filepath.Join(tdir, "hello_world"))
	})
	r.NoError(err)

}

func Test_NewCmd_API(t *testing.T) {
	if envy.Get("GO111MODULE", "off") == "on" {
		t.Skip("CURRENTLY NOT SUPPORTED")
	}
	args := []string{
		"new",
		"hello_world",
		"--skip-pop",
		"--api",
		"--vcs=none",
	}
	r := require.New(t)
	err := call(args, func(tdir string) {
		r.DirExists(filepath.Join(tdir, "hello_world"))
	})
	r.NoError(err)

}

func Test_NewCmd_WithPopSQLite3(t *testing.T) {
	if envy.Get("GO111MODULE", "off") == "on" {
		t.Skip("CURRENTLY NOT SUPPORTED")
	}
	r := require.New(t)

	args := []string{
		"new",
		"hello_world",
		"--db-type=sqlite3",
		"--skip-webpack",
		"--vcs=none",
	}

	err := call(args, func(tdir string) {
		r.DirExists(filepath.Join(tdir, "hello_world"))
		r.FileExists(filepath.Join(tdir, "hello_world", "database.yml"))
	})
	r.NoError(err)

}
