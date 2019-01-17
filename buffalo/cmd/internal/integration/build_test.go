// +build integration_test

package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Build_Nominal(t *testing.T) {
	r := require.New(t)
	args := []string{
		"new",
		"build_nominal",
		"--skip-pop",
		"--skip-webpack",
		"--vcs=none",
	}
	err := call(args, func(tdir string) {
		fmt.Println("### tdir ->", tdir)
		ad := filepath.Join(tdir, "build_nominal")
		fmt.Println("### ad ->", ad)
		r.DirExists(ad)
		os.Chdir(ad)
		pwd, _ := os.Getwd()
		fmt.Println("### pwd ->", pwd)

		args = []string{"build"}
		err := exec(args)
		r.NoError(err)
	})
	r.NoError(err)

}
