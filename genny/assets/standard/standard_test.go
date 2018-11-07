package standard

import (
	"context"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(nil)
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)

	files := []string{
		"public/assets/application.css",
		"public/assets/application.js",
		"public/assets/buffalo.css",
		"public/assets/images/favicon.ico",
	}
	r.Len(res.Files, len(files))
	for i, f := range res.Files {
		r.Equal(files[i], f.Name())
	}
}
