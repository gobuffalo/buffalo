package grift

import (
	"testing"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	err := run.WithNew(New(&Options{
		Args: []string{"foo"},
	}))
	r.NoError(err)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal("grifts/foo.go", f.Name())
	body := f.String()
	r.Contains(body, `var _ = Add("foo", func(c *Context) error`)
}

func Test_New_Namespaced(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	err := run.WithNew(New(&Options{
		Args: []string{"foo:bar"},
	}))
	r.NoError(err)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal("grifts/bar.go", f.Name())
	body := f.String()
	r.Contains(body, `Add("bar", func(c *Context) error`)
}

func Test_New_No_Name(t *testing.T) {
	r := require.New(t)

	_, err := New(&Options{})
	r.Error(err)
}
