package actions

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	packr "github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

func runner() *genny.Runner {
	run := gentest.NewRunner()
	run.Disk.AddBox(packr.New("actions/start/test", "../actions/_fixtures/start"))
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

	box := packr.New("actions/normal/test", "../actions/_fixtures/normal")
	appGo, err := box.FindString("actions/app.go")
	r.NoError(err)

	userGo, err := box.FindString("actions/user.go")
	r.NoError(err)
	f, err := res.Find("actions/user.go")
	r.NoError(err)
	r.Equal(strings.TrimSpace(userGo), strings.TrimSpace(f.String()))

	f, err = res.Find("actions/app.go")
	r.NoError(err)
	r.Equal(strings.TrimSpace(appGo), strings.TrimSpace(f.String()))

	ind, err := box.FindString("templates/user/index.html")
	r.NoError(err)
	f, err = res.Find("templates/user/index.html")
	r.NoError(err)
	r.Equal(strings.TrimSpace(ind), strings.TrimSpace(f.String()))

	tst, err := box.FindString("actions/user_test.go.tmpl")
	r.NoError(err)

	f, err = res.Find("actions/user_test.go")
	r.NoError(err)
	r.Equal(strings.TrimSpace(tst), strings.TrimSpace(f.String()))
}
