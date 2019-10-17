package info

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/clara/genny/rx"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/packd"
	"github.com/stretchr/testify/require"
)

func Test_pkgChecks(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	run := gentest.NewRunner()

	opts := &Options{
		App: meta.New("."),
		Out: rx.NewWriter(bb),
	}

	box := packd.NewMemoryBox()
	box.AddString("go.mod", "module foo")
	run.WithRun(pkgChecks(opts, box))

	r.NoError(run.Run())

	res := bb.String()
	r.Contains(res, "Buffalo: go.mod")
}
