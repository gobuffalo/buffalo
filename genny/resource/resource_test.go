package resource

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/genny/movinglater/attrs"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

type pass struct {
	Name    string
	Options Options
}

func Test_New(t *testing.T) {
	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	if err != nil {
		t.Fatal(err)
	}
	table := []pass{
		{"default", Options{Name: "widget", Attrs: ats}},
		// {"nested", Options{Name: "admin/widget", Attrs: ats}},
		// {"deep_nested", Options{Name: "depp/admin/widget", Attrs: ats}},
		// {"skip_migration", Options{Name: "widget", Attrs: ats, SkipMigration: true}},
		// {"skip_model", Options{Name: "widget", Attrs: ats, SkipModel: true}},
		// {"use_model", Options{Name: "widget", Attrs: ats, UseModel: true, Model: "gadget"}},
	}

	for _, tt := range table {
		t.Run(tt.Name, func(st *testing.T) {
			r := require.New(st)
			g, err := New(&tt.Options)
			r.NoError(err)

			run := gentest.NewRunner()
			run.Disk.AddBox(packr.New("./_fixtures/coke", "./_fixtures/coke"))
			g.Transformer(genny.NewTransformer(".", func(f genny.File) (genny.File, error) {
				s := strings.Replace(f.String(), "\t", "  ", -1)
				return genny.NewFileS(f.Name(), s), nil
			}))
			run.With(g)

			r.NoError(run.Run())

			res := run.Results()

			r.Len(res.Commands, 1)

			c := res.Commands[0]
			r.Equal("buffalo-pop pop g model widget desc:nulls.Text", strings.Join(c.Args, " "))

			for _, f := range res.Files {
				fmt.Println(f.Name())
			}

			exp := packr.New(tt.Name, filepath.Join("_fixtures", tt.Name))
			r.Len(res.Files, len(exp.List()))

			for _, n := range exp.List() {
				f, err := res.Find(strings.TrimSuffix(n, ".tmpl"))
				r.NoError(err)
				s, err := exp.FindString(n)
				r.NoError(err)
				r.Equal(strings.TrimSpace(s), strings.TrimSpace(f.String()))
			}

		})
	}
}
