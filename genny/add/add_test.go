package add

import (
	"bytes"
	"fmt"
	"go/build"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/plugins/plugdeps"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		App: meta.New("."),
		Plugins: []plugdeps.Plugin{
			{Binary: "buffalo-pop", GoGet: "github.com/gobuffalo/buffalo-pop"},
			{Binary: "buffalo-hello.rb", Local: "./plugins/buffalo-hello.rb"},
		},
	})
	r.NoError(err)

	run := gentest.NewRunner()
	c := build.Default
	run.Disk.Add(genny.NewFile(filepath.Join(c.GOPATH, "bin", "buffalo-pop"), &bytes.Buffer{}))

	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	ecmds := []string{}
	r.Len(res.Commands, len(ecmds))
	for i, c := range res.Commands {
		r.Equal(ecmds[i], strings.Join(c.Args, " "))
	}

	efiles := []string{"bin/buffalo-pop", "config/buffalo-plugins.toml"}
	r.Len(res.Files, len(efiles))

	for i, f := range res.Files {
		r.True(strings.HasSuffix(f.Name(), efiles[i]), fmt.Sprintf("Not found %v", efiles[i]))
	}
}
