package info

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/clara/genny/rx"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_appDetails(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	app := meta.New(".")
	app.Bin = "paris elephant chevrolet"

	bb := &bytes.Buffer{}

	opts := &Options{
		App: app,
		Out: rx.NewWriter(bb),
	}

	run.WithRun(appDetails(opts))

	r.NoError(run.Run())

	r.Contains(bb.String(), "paris elephant chevrolet")
}
