package meta

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ModulesPackageName(t *testing.T) {
	r := require.New(t)
	tmp := os.TempDir()
	modsOn = true

	r.NoError(os.Chdir(tmp))
	r.NoError(ioutil.WriteFile("go.mod", []byte("module github.com/wawandco/zekito"), 0644))

	a := New(tmp)
	r.Equal("github.com/wawandco/zekito", a.PackagePkg)
}
