package resource

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/genny/movinglater/attrs"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

var coke = packr.New("coke", "./_fixtures")

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
			run.Root = filepath.Join("_fixtures", "core")
			run.With(g)

			r.NoError(run.Run())

			res := run.Results()

			r.Len(res.Commands, 1)

			c := res.Commands[0]
			r.Equal("buffalo-pop pop g model widget desc:nulls.Text", strings.Join(c.Args, " "))

			box := packr.New(tt.Name, filepath.Join("_fixtures", tt.Name))

			for _, f := range res.Files {
				fmt.Println(f.Name())
			}

			r.Len(res.Files, len(box.List()))

			for _, n := range box.List() {
				f, err := res.Find(n)
				r.NoError(err)
				fmt.Println("### f.Name() ->", f.Name())
			}

		})
	}
}
