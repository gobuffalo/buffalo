package build

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

// TODO: once `buffalo new` is converted to use genny
// create an integration test that first generates a new application
// and then tries to build using genny/build.
var coke = packr.NewBox("../build/_fixtures/coke")

var cokeRunner = func() *genny.Runner {
	run := gentest.NewRunner()
	run.Disk.AddBox(coke)
	run.Root = coke.Path
	return run
}

func Test_New(t *testing.T) {
	r := require.New(t)

	run := cokeRunner()

	opts := &Options{
		WithAssets:  true,
		Environment: "bar",
		App:         meta.New("."),
	}
	opts.App.Bin = "bin/foo"
	r.NoError(run.WithNew(New(opts)))
	run.Root = opts.App.Root

	r.NoError(run.Run())

	res := run.Results()

	// we should never leave any files modified or dropped
	r.Len(res.Files, 0)

	cmds := []string{"go get ./...", "go build -tags bar -o bin/foo"}
	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}
}
