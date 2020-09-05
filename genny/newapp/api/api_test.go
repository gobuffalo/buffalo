package api

import (
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/genny/newapp/core"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func init() {
	// normalize command output
	envy.Set("GO_BIN", "go")
}

func Test_New(t *testing.T) {
	r := require.New(t)

	app := meta.Named("api", ".")
	(&app).PackageRoot("api")
	app.WithModules = false
	app.AsAPI = true
	app.AsWeb = false

	gg, err := New(&Options{
		Options: &core.Options{
			App: app,
		},
	})
	r.NoError(err)

	run := gentest.NewRunner()
	run.WithGroup(gg)

	r.NoError(run.Run())

	res := run.Results()

	cmds := []string{
		"go mod init api",
	}
	r.Len(res.Commands, len(cmds))

	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}

	expected := commonExpected
	for _, e := range expected {
		_, err = res.Find(e)
		r.NoError(err)
	}

	f, err := res.Find("actions/render.go")
	r.NoError(err)
	r.Contains(f.String(), `DefaultContentType: "application/json",`)

	f, err = res.Find("actions/home.go")
	r.NoError(err)
	r.Contains(f.String(), `return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "Welcome to Buffalo!"}))`)

	f, err = res.Find("actions/app.go")
	r.NoError(err)
	r.Contains(f.String(), `i18n "github.com/gobuffalo/mw-i18n"`)
	r.Contains(f.String(), `var T *i18n.Translator`)
	r.Contains(f.String(), `func translations() buffalo.MiddlewareFunc {`)

	f, err = res.Find("locales/all.en-us.yaml")
	r.NoError(err)
	r.Contains(f.String(), `translation: "Welcome to Buffalo (EN)"`)

	unexpected := []string{
		"Dockerfile",
		"database.yml",
		"models/models.go",
		"go.mod",
		".buffalo.dev.yml",
		"assets/css/application.scss.css",
		"public/assets/application.js",
	}

	for _, u := range unexpected {
		_, err = res.Find(u)
		r.Error(err)
	}
}

var commonExpected = []string{
	"main.go",
	"actions/app.go",
	"actions/actions_test.go",
	"actions/render.go",
	"actions/home.go",
	"actions/home_test.go",
	"fixtures/sample.toml",
	"grifts/init.go",
	".codeclimate.yml",
	".env",
	"inflections.json",
	"README.md",
}
