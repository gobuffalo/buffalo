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
	os.Chdir(dirname)
	_, err = os.Create(".buffalo.dev.yml")
	projectDir, err := os.Getwd()
	r.NoError(err)
	for i := 0; i < 4; i++ {
		dirname, err := ioutil.TempDir(dirname, "nested_dirs")
		r.NoError(err)
		r.NoError(os.Chdir(dirname))
		dirname, err = os.Getwd()
		r.NoError(err)
	}
	r.True(cmd.SwitchToBuffaloProjectDir())
	dirname, err = os.Getwd()
	r.Equal(dirname, projectDir)
	os.Chdir(os.TempDir())
	os.RemoveAll(dirnameRemove)
}
