package grift

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Args: []string{"foo"},
	})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal(filepath.Join("grifts", "foo.go"), f.Name())
	body := f.String()
	r.Contains(body, `var _ = Add("foo", func(c *Context) error`)
}

func Test_New_Namespaced(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Args: []string{"foo:bar"},
	})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal(filepath.Join("grifts", "bar.go"), f.Name())
	body := f.String()
	r.Contains(body, `Add("bar", func(c *Context) error`)
}

func Test_New_No_Name(t *testing.T) {
	r := require.New(t)

	_, err := New(&Options{})
	r.Error(err)
}
