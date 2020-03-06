package info

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/clara/genny/rx"

	"github.com/gobuffalo/packd"

	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_configs(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	bb := &bytes.Buffer{}

	app := meta.New(".")
	opts := &Options{
		App: app,
		Out: rx.NewWriter(bb),
	}

	box := packd.NewMemoryBox()
	box.AddString("buffalo-app.toml", "app")
	box.AddString("buffalo-plugins.toml", "plugins")
	run.WithRun(configs(opts, box))

	r.NoError(run.Run())

	x := bb.String()
	r.Contains(x, "Buffalo: config/buffalo-app.toml\napp")
	r.Contains(x, "Buffalo: config/buffalo-plugins.toml\nplugins")
}
