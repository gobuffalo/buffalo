package resource

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/flect/name"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/attrs"

	"github.com/stretchr/testify/require"
)

func Test_Resource_New(t *testing.T) {
	r := require.New(t)

	gg, err := New(&Options{
		Attrs: attrs.NamedAttrs{
			Name: name.New("widget"),
		},
	})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.Disk.Add(genny.NewFile(filepath.Join("actions", "app.go"), strings.NewReader(actionsApp)))
	gg.With(run)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Contains(f.String(), `app.Resource("/widgets", WidgetsResource{})`)

}

const actionsApp = `
package actions

import "github.com/gobuffalo/buffalo"

func App() *buffalo.App {
	var app *buffalo.app
	if app == nil {

	}
	return app
}
`
