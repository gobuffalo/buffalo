package resource

import (
	"path/filepath"
	"testing"

	"github.com/gobuffalo/genny/gentest"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

var coke = packr.New("coke", "./_fixtures")

type pass struct {
	Name  string
	Args  []string
	Flags []string
}

func Test_New(t *testing.T) {

	table := []pass{
		{"default", []string{"widget", "name", "desc:nulls.Text"}, []string{}},
		{"nested", []string{"admin/widget", "name", "desc:nulls.Text"}, []string{}},
		{"deep_nested", []string{"deep/admin/widget", "name", "desc:nulls.Text"}, []string{}},
		{"skip_migration", []string{"widget", "name", "desc:nulls.Text"}, []string{"--skip-migration"}},
		{"skip_model", []string{"widget", "name", "desc:nulls.Text"}, []string{"--skip-model"}},
		{"use_model", []string{"widget", "name", "desc:nulls.Text"}, []string{"--use-model", "gadget"}},
	}

	for _, tt := range table {
		t.Run(tt.Name, func(st *testing.T) {
			r := require.New(st)
			g, err := New(&Options{})
			r.NoError(err)

			run := gentest.NewRunner()
			run.With(g)

			r.NoError(run.Run())

			res := run.Results()

			r.Len(res.Commands, 0)

			box := packr.New(tt.Name, filepath.Join("_fixtures", tt.Name))
			r.Len(res.Files, len(box.List()))

			f := res.Files[0]
			r.Equal("example.txt", f.Name())

		})
	}
}
