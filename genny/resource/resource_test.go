package resource

import (
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/flect/name"
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

			r.Len(res.Files, 9)

			nn := name.New(tt.Options.Name).Pluralize().String()
			for _, s := range []string{"_form", "edit", "index", "new", "show"} {
				p := path.Join("templates", nn, s+".html")
				_, err = res.Find(p)
				r.NoError(err)
			}

			exp := packr.New(tt.Name, filepath.Join("_fixtures", tt.Name))
			gentest.CompareFiles(exp.List(), res.Files)

			for _, n := range exp.List() {
				f, err := res.Find(strings.TrimSuffix(n, ".tmpl"))
				r.NoError(err)
				s, err := exp.FindString(n)
				r.NoError(err)

				clean := func(s string) string {
					s = strings.TrimSpace(s)
					return s
				}
				r.Equal(clean(s), clean(f.String()))
			}

		})
	}
}

func Test_New_SkipTemplates(t *testing.T) {
	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	if err != nil {
		t.Fatal(err)
	}
	table := []pass{
		{"default", Options{Name: "widget", Attrs: ats}},
		{"nested", Options{Name: "admin/widget", Attrs: ats}},
	}

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")

	for _, tt := range table {
		t.Run(tt.Name, func(st *testing.T) {
			tt.Options.App = app
			tt.Options.SkipTemplates = true
			r := require.New(st)
			g, err := New(&tt.Options)
			r.NoError(err)

			run := runner()
			run.With(g)
			r.NoError(run.Run())

			res := run.Results()

			r.Len(res.Commands, 1)

			nn := name.New(tt.Options.Name).Pluralize().String()
			for _, s := range []string{"_form", "edit", "index", "new", "show"} {
				p := path.Join("templates", nn, s+".html")
				_, err = res.Find(p)
				r.Error(err)
			}

			r.Len(res.Files, 3)
		})
	}
}

func Test_New_API(t *testing.T) {
	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	if err != nil {
		t.Fatal(err)
	}
	table := []pass{
		{"default", Options{Name: "widget", Attrs: ats}},
		{"nested", Options{Name: "admin/widget", Attrs: ats}},
	}

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")
	app.AsAPI = true

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

			nn := name.New(tt.Options.Name).Pluralize().String()
			for _, s := range []string{"_form", "edit", "index", "new", "show"} {
				p := path.Join("templates", nn, s+".html")
				_, err = res.Find(p)
				r.Error(err)
			}

			r.Len(res.Files, 3)
		})
	}
}

func Test_New_UseModel(t *testing.T) {
	r := require.New(t)

	ats, err := attrs.ParseArgs("name", "desc:nulls.Text")
	r.NoError(err)

	app := meta.New(".")
	app.PackageRoot("github.com/markbates/coke")

	opts := &Options{
		App:   app,
		Name:  "Widget",
		Model: "User",
		Attrs: ats,
	}
	g, err := New(opts)
	r.NoError(err)

	run := runner()
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 1)

	c := res.Commands[0]
	r.Equal("buffalo-pop pop g model user desc:nulls.Text", strings.Join(c.Args, " "))

	r.Len(res.Files, 9)

	for _, s := range []string{"_form", "edit", "index", "new", "show"} {
		p := path.Join("templates", "widgets", s+".html")
		_, err = res.Find(p)
		r.NoError(err)
	}

	f, err := res.Find("actions/widgets.go")
	r.NoError(err)
	r.Contains(f.String(), "users := &models.Users{}")

}
