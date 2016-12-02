package cmd

var newTemplates = map[string]string{
	"main.go":                    nMain,
	"refresh.yml":                nRefresh,
	"actions/app.go":             nApp,
	"actions/home.go":            nHomeHandler,
	"actions/home_test.go":       nHomeHandlerTest,
	"actions/render.go":          nRender,
	"grifts/routes.go":           nGriftRoutes,
	"templates/index.html":       nIndexHTML,
	"templates/application.html": nApplicationHTML,
	"assets/application.js":      "",
	"assets/application.css":     nApplicationCSS,
	".gitignore":                 nGitignore,
}

var nMain = `package main

import (
	"log"
	"net/http"

	"{{.actionsPath}}"
)

func main() {
	log.Fatal(http.ListenAndServe(":3000", actions.App()))
}

`
var nApp = `package actions

import (
	"net/http"

	"github.com/markbates/buffalo"
)

func App() http.Handler {
	a := buffalo.Automatic(buffalo.Options{})
	a.Env = "development"

	a.ServeFiles("/assets", assetsPath())
	a.GET("/", HomeHandler)

	return a
}
`

var nRender = `package actions

import (
	"net/http"
	"path"
	"runtime"

	"github.com/markbates/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(&render.Options{
		TemplatesPath: fromHere("../templates"),
		HTMLLayout:    "application.html",
	})
}

func assetsPath() http.Dir {
	return http.Dir(fromHere("../assets"))
}

func fromHere(p string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), p)
}
`

var nHomeHandler = `package actions

import "github.com/markbates/buffalo"

func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}
`

var nHomeHandlerTest = `package actions_test

import (
	"testing"

	"{{.actionsPath}}"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_HomeHandler(t *testing.T) {
	r := require.New(t)

	w := willie.New(actions.App())
	res := w.Request("/").Get()

	r.Equal(200, res.Code)
	r.Contains(res.Body.String(), "Welcome to Buffalo!")
}
`

var nIndexHTML = `<h1>Welcome to Buffalo!</h1>`

var nApplicationHTML = `<html>
<head>
  <meta charset="utf-8">
  <title>Buffalo - {{ .name }}</title>
  <link rel="stylesheet" href="/assets/application.css" type="text/css" media="all" />
</head>
<body>
  {{"{{"}} yield {{"}}"}}

  <script src="/assets/application.js" type="text/javascript" charset="utf-8"></script>
</body>
</html>
`

var nApplicationCSS = `body {
  font-family: helvetica;
}
`

var nGitignore = `vendor/
**/*.log
**/*.sqlite
bin/
node_modules/
{{ .name }}
`

var nGriftRoutes = `package grifts

import (
	"os"

	"github.com/markbates/buffalo"
	. "github.com/markbates/grift/grift"
	"{{.actionsPath}}"
	"github.com/olekukonko/tablewriter"
)

var _ = Add("routes", func(c *Context) error {
	a := actions.App().(*buffalo.App)
	routes := a.Routes()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Path", "Handler"})
	for _, r := range routes {
		table.Append([]string{r.Method, r.Path, r.HandlerName})
	}
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return nil
})`

var nRefresh = `app_root: .
ignored_folders:
- vendor
- log
- tmp
included_extensions:
- .go
- .html
build_path: /tmp
build_delay: 200ns
binary_name: {{.name}}-build
command_flags: []
enable_colors: true`
