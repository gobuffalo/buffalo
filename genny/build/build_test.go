package build

import (
	"bytes"
	"io"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

var coke = packr.NewBox("../build/_fixtures/coke")

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	coke.Walk(func(path string, file packr.File) error {
		bb := &bytes.Buffer{}
		io.Copy(bb, file)
		f := genny.NewFile(path, bb)
		run.Disk.Add(f)
		return nil
	})
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 0)
}
