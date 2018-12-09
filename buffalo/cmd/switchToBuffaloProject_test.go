package cmd

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gobuffalo/buffalo/buffalo/cmd"
	"github.com/stretchr/testify/require"
)

func Test_SwitchToBuffaloProject(t *testing.T) {
	r := require.New(t)
	r.NoError(os.Chdir(os.TempDir()))
	dirname, err := os.Getwd()
	r.NoError(err)
	dirname, err = ioutil.TempDir(dirname, "buffalo_testing_dir")
	r.NoError(err)
	dirnameRemove := dirname
	r.NoError(os.Chdir(dirname))
	_, err = os.Create(".buffalo.dev.yml")
	r.NoError(err)
	projectDir, err := os.Getwd()
	r.NoError(err)
	for i := 0; i < 4; i++ {
		dirname, err = os.Getwd()
		r.NoError(err)
		dirname, err = ioutil.TempDir(dirname, "nested_dirs")
		r.NoError(err)
		r.NoError(os.Chdir(dirname))
	}
	r.True(cmd.SwitchToBuffaloProjectDir())
	dirname, err = os.Getwd()
	r.NoError(err)
	r.Equal(dirname, projectDir)
	r.NoError(os.Chdir(os.TempDir()))
	r.NoError(os.RemoveAll(dirnameRemove))
}
