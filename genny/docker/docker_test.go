package docker

import (
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)

	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal(".dockerignore", f.Name())

	f = res.Files[1]
	r.Equal("Dockerfile", f.Name())
	r.Contains(f.String(), "multi-stage")
}

func Test_New_Standard(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Style: "standard",
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)

	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal(".dockerignore", f.Name())

	f = res.Files[1]
	r.Equal("Dockerfile", f.Name())
	r.NotContains(f.String(), "multi-stage")
}
