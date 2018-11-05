package newapp

import (
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	gg, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	run.WithGroup(gg)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 4)
	r.Len(res.Files, 22)
}
