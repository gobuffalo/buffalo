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
	r.Len(res.Files, 3)

	f := res.Files[0]
	r.Equal("public/assets/application.css", f.Name())

	f = res.Files[1]
	r.Equal("public/assets/application.js", f.Name())

	f = res.Files[2]
	r.Equal("public/assets/images/favicon.ico", f.Name())
}
