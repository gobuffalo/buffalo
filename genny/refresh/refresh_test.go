package refresh

import (
	"context"
	"testing"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	app := meta.New(".")
	app.Name = name.New("foo")
	g, err := New(&Options{app})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Contains(f.String(), "binary_name: foo-build")
}
