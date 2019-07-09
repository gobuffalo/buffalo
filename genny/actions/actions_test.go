package actions

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func compare(a, b string) bool {
	a = strings.TrimSpace(a)
	a = strings.Replace(a, "\r", "", -1)
	b = strings.TrimSpace(b)
	b = strings.Replace(b, "\r", "", -1)
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	res := cmp.Equal(a, b)
	if !res {
		fmt.Println(cmp.Diff(a, b))
	}
	return res
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

func Test_New_Multi_Existing(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:    "user",
		Actions: []string{"show", "edit"},
	})
	r.NoError(err)

	run := gentest.NewRunner()
	ins := packr.New("Test_New_Multi_Existing_input", "../actions/_fixtures/inputs/existing")
	for _, n := range ins.List() {
		x, err := ins.FindString(n)
		r.NoError(err)
		n = strings.TrimSuffix(n, ".tmpl")
		run.Disk.Add(genny.NewFileS(n, x))
	}
	run.With(g)

	err = run.Run()
	r.NoError(err)

	res := run.Results()

	r.Len(res.Commands, 0)

	box := packr.New("genny/actions/Test_New_Multi_Existing", "../actions/_fixtures/outputs/existing")

	files := []string{"actions/user.go.tmpl", "actions/app.go.tmpl", "actions/user_test.go.tmpl", "templates/user/show.html", "templates/user/edit.html"}

	for _, s := range files {
		x, err := box.FindString(s)
		r.NoError(err)
		f, err := res.Find(strings.TrimSuffix(s, ".tmpl"))
		r.NoError(err)
		r.True(compare(x, f.String()))
	}
}

func Test_New_SkipTemplates(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:          "user",
		Actions:       []string{"index"},
		SkipTemplates: true,
	})
	r.NoError(err)

	run := runner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)

	files := []string{"templates/user/index.html"}

	for _, s := range files {
		_, err := res.Find(s)
		r.Error(err)
	}
}
