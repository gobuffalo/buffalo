package meta

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/require"
)

func Test_ModulesPackageName(t *testing.T) {
	r := require.New(t)
	tmp := os.TempDir()
	modsOn = true

	r.NoError(os.Chdir(tmp))

	tcases := []struct {
		Content     string
		PackageName string
	}{
		{"module github.com/wawandco/zekito", "github.com/wawandco/zekito"},
		{"module zekito", "zekito"},
		{"module gopkg.in/some/friday.v2", "gopkg.in/some/friday.v2"},
		{"", "zekito"},
	}

	for _, tcase := range tcases {
		envy.Set("GOPATH", tmp)

		t.Run(tcase.Content, func(st *testing.T) {
			r := require.New(st)

			r.NoError(ioutil.WriteFile("go.mod", []byte(tcase.Content), 0644))

			a := New(filepath.Join(tmp, "zekito"))
			r.Equal(tcase.PackageName, a.PackagePkg)
		})
	}
}
