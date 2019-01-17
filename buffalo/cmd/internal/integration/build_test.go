// +build integration_test

package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test_New_Build_Nominal creates a new nominal
// app and then builds it
func Test_New_Build_Nominal(t *testing.T) {
	r := require.New(t)
	args := []string{
		"new",
		"build_nominal",
		"--skip-pop",
		"--skip-webpack",
		"--vcs=none",
	}
	err := call(args, func(tdir string) {
		ad := filepath.Join(tdir, "build_nominal")
		r.DirExists(ad)
		os.Chdir(ad)

		args = []string{"build"}
		err := exec(args)
		r.NoError(err)
	})
	r.NoError(err)

}

// Test_New_Build_Api creates a new API
// app and then builds it
func Test_New_Build_Api(t *testing.T) {
	r := require.New(t)
	args := []string{
		"new",
		"build_api",
		"--skip-pop",
		"--api",
		"--vcs=none",
	}
	err := call(args, func(tdir string) {
		ad := filepath.Join(tdir, "build_api")
		r.DirExists(ad)
		os.Chdir(ad)

		args = []string{"build"}
		err := exec(args)
		r.NoError(err)
	})
	r.NoError(err)

}
func Test_New_Build_Sqlite(t *testing.T) {
	r := require.New(t)

	args := []string{
		"new",
		"build_sqlite",
		"--db-type=sqlite3",
		"--skip-webpack",
		"--vcs=none",
	}

	err := call(args, func(tdir string) {
		ad := filepath.Join(tdir, "build_sqlite")
		r.DirExists(ad)
		r.FileExists(filepath.Join(ad, "database.yml"))
		os.Chdir(ad)

		args = []string{"build"}
		err := exec(args)
		r.NoError(err)
	})
	r.NoError(err)

}
