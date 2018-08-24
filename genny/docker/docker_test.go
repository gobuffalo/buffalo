package docker

import (
	"context"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)
	r.NotNil(g)

	run := genny.DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal(".dockerignore", f.Name())

	f = res.Files[1]
	r.Equal("Dockerfile", f.Name())
	r.Contains(f.String(), "multi-stage Dockerfile")
}

func Test_New_standard(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Style: "standard",
	})
	r.NoError(err)
	r.NotNil(g)

	run := genny.DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 2)

	f := res.Files[0]
	r.Equal(".dockerignore", f.Name())

	f = res.Files[1]
	r.Equal("Dockerfile", f.Name())
	r.NotContains(f.String(), "multi-stage Dockerfile")
}

func Test_New_none(t *testing.T) {
	r := require.New(t)

	_, err := New(&Options{
		Style: "none",
	})
	r.Error(err)
}
