package resource

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/genny/movinglater/attrs"
	"github.com/gobuffalo/meta"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

type pass struct {
	Name    string
	Options Options
}

func runner() *genny.Runner {
	run := gentest.NewRunner()
	box := packr.New("./_fixtures/coke", "./_fixtures/coke")
	box.Walk(func(path string, file packr.File) error {
		path = strings.TrimSuffix(path, ".tmpl")
		run.Disk.Add(genny.NewFile(path, file))
		return nil
	})
	return run
}

func Test_New(t *testing.T) {
	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	if err != nil {
		t.Fatal(err)
	}
	table := []pass{
		{"default", Options{Name: "widget", Attrs: ats}},
		{"nested", Options{Name: "admin/widget", Attrs: ats}},
		// {"deep_nested", Options{Name: "depp/admin/widget", Attrs: ats}},
		// {"skip_migration", Options{Name: "widget", Attrs: ats, SkipMigration: true}},
		// {"skip_model", Options{Name: "widget", Attrs: ats, SkipModel: true}},
		// {"use_model", Options{Name: "widget", Attrs: ats, UseModel: true, Model: "gadget"}},
	}

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")
	for _, tt := range table {
		t.Run(tt.Name, func(st *testing.T) {
			tt.Options.App = app
			r := require.New(st)
			g, err := New(&tt.Options)
			r.NoError(err)

			run := runner()
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
			gentest.CompareFiles(exp.List(), res.Files)

			for _, n := range exp.List() {
				f, err := res.Find(strings.TrimSuffix(n, ".tmpl"))
				r.NoError(err)
				s, err := exp.FindString(n)
				r.NoError(err)

				clean := func(s string) string {
					// fmt.Println(s)
					s = strings.TrimSpace(s)
					// s = strings.Replace(s, "\n", "", -1)
					// s = strings.Replace(s, "\t", "", -1)
					return s
				}
				r.Equal(clean(s), clean(f.String()))
			}

		})
	}
}
