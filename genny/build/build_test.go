package build

import (
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

// TODO: once `buffalo new` is converted to use genny
// create an integration test that first generates a new application
// and then tries to build using genny/build.
var coke = packr.New("github.com/gobuffalo/buffalo/genny/build/build_test", "../build/_fixtures/coke")

var cokeRunner = func() *genny.Runner {
	run := gentest.NewRunner()
	run.Disk.AddBox(coke)
	run.Root = coke.Path
	return run
}

var eq = func(r *require.Assertions, s string, c *exec.Cmd) {
	if runtime.GOOS == "windows" {
		s = strings.Replace(s, "bin/build", `bin\build.exe`, 1)
		s = strings.Replace(s, "bin/foo", `bin\foo.exe`, 1)
	}
	r.Equal(s, strings.Join(c.Args, " "))
}

func Test_New(t *testing.T) {
	envy.Temp(func() {
		envy.Set(envy.GO111MODULE, "off")
		r := require.New(t)

		run := cokeRunner()

		opts := &Options{
			WithAssets:    true,
			WithBuildDeps: true,
			Environment:   "bar",
			App:           meta.New("."),
		}
		opts.App.Bin = "bin/foo"
		r.NoError(run.WithNew(New(opts)))
		run.Root = opts.App.Root

		r.NoError(run.Run())

		res := run.Results()

		// we should never leave any files modified or dropped
		r.Len(res.Files, 0)

		cmds := []string{"go get -tags bar ./...", "go build -tags bar -o bin/foo"}
		r.Len(res.Commands, len(cmds))
		for i, c := range res.Commands {
			eq(r, cmds[i], c)
		}
	})
}

func Test_NewWithoutBuildDeps(t *testing.T) {
	envy.Temp(func() {
		envy.Set(envy.GO111MODULE, "off")
		r := require.New(t)

		run := cokeRunner()

		opts := &Options{
			WithAssets:    false,
			WithBuildDeps: false,
			Environment:   "bar",
			App:           meta.New("."),
		}
		opts.App.Bin = "bin/foo"
		r.NoError(run.WithNew(New(opts)))
		run.Root = opts.App.Root

		r.NoError(run.Run())

		res := run.Results()

		cmds := []string{"go build -tags bar -o bin/foo"}
		r.Len(res.Commands, len(cmds))
		for i, c := range res.Commands {
			eq(r, cmds[i], c)
		}
	})
}
