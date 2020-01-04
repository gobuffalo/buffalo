package web

import (
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/genny/newapp/core"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func init() {
	// normalize command output
	envy.Set("GO_BIN", "go")
}

func Test_New(t *testing.T) {
	r := require.New(t)

	app := meta.New(".")
	app.WithModules = true
	envy.Set(envy.GO111MODULE, "on")

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
		"go mod init web",
		"go get github.com/gobuffalo/buffalo@v0.15.3",
		"go get ./...",
		"go mod tidy",
	}
	r.Len(res.Commands, len(cmds))

	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}

	for _, e := range commonExpected {
		_, err = res.Find(e)
		r.NoError(err)
	}

	f, err := res.Find("actions/render.go")
	r.NoError(err)

	body := f.String()
	r.Contains(body, `TemplatesBox: packr.New("app:templates", "../templates"),`)
	r.NotContains(body, `DefaultContentType: "application/json",`)
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
	"locales/all.en-us.yaml",
	"public/robots.txt",
	"templates/_flash.plush.html",
	"templates/application.plush.html",
	"templates/index.plush.html",
}
