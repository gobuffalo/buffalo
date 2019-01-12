package actions

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

func compare(a, b string) bool {
	a = strings.TrimSpace(a)
	a = strings.Replace(a, "\r", "", -1)
	b = strings.TrimSpace(b)
	b = strings.Replace(b, "\r", "", -1)
	return a == b
}

func runner() *genny.Runner {
	run := gentest.NewRunner()
	run.Disk.AddBox(packr.New("actions/start/test", "../actions/_fixtures/inputs/clean"))
	return run
}

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:    "user",
		Actions: []string{"index"},
	})
	r.NoError(err)

	run := runner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	// r.Len(res.Files, 4)

	box := packr.New("genny/actions/Test_New", "../actions/_fixtures/outputs/clean")

	files := []string{"actions/user.go.tmpl", "actions/app.go.tmpl", "actions/user_test.go.tmpl", "templates/user/index.html"}

	for _, s := range files {
		x, err := box.FindString(s)
		r.NoError(err)
		f, err := res.Find(strings.TrimSuffix(s, ".tmpl"))
		r.NoError(err)
		r.True(compare(x, f.String()))
	}
}

func Test_New_Multi(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:    "user",
		Actions: []string{"show", "edit"},
	})
	r.NoError(err)

	run := runner()
	run.With(g)

	err = run.Run()
	r.NoError(err)

	res := run.Results()

	r.Len(res.Commands, 0)

	box := packr.New("genny/actions/Test_New_Multi", "../actions/_fixtures/outputs/multi")

	files := []string{"actions/user.go.tmpl", "actions/app.go.tmpl", "actions/user_test.go.tmpl", "templates/user/show.html", "templates/user/edit.html"}

	for _, s := range files {
		x, err := box.FindString(s)
		r.NoError(err)
		f, err := res.Find(strings.TrimSuffix(s, ".tmpl"))
		r.NoError(err)
		r.True(compare(x, f.String()))
	}
}
