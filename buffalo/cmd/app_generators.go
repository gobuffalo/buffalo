package cmd

import (
	"os/exec"

	"github.com/markbates/gentronics"
)

func newAppGenerator() *gentronics.Generator {
	g := gentronics.New()
	g.Add(gentronics.NewFile("main.go", nMain))
	g.Add(gentronics.NewFile(".buffalo.dev.yml", nRefresh))
	g.Add(gentronics.NewFile("actions/app.go", nApp))
	g.Add(gentronics.NewFile("actions/home.go", nHomeHandler))
	g.Add(gentronics.NewFile("actions/home_test.go", nHomeHandlerTest))
	g.Add(gentronics.NewFile("actions/render.go", nRender))
	g.Add(gentronics.NewFile("grifts/routes.go", nGriftRoutes))
	g.Add(gentronics.NewFile("templates/index.html", nIndexHTML))
	g.Add(gentronics.NewFile("templates/application.html", nApplicationHTML))
	g.Add(gentronics.NewFile("assets/application.js", ""))
	g.Add(gentronics.NewFile("assets/application.css", nApplicationCSS))
	g.Add(gentronics.NewFile(".gitignore", nGitignore))
	g.Add(gentronics.NewCommand(goGet("github.com/markbates/refresh/...")))
	g.Add(gentronics.NewCommand(goInstall("github.com/markbates/refresh")))
	g.Add(gentronics.NewCommand(goGet("github.com/markbates/grift/...")))
	g.Add(gentronics.NewCommand(goInstall("github.com/markbates/grift")))
	g.Add(newJQueryGenerator())
	g.Add(newSodaGenerator())
	g.Add(gentronics.NewCommand(appGoGet()))
	return g
}

func appGoGet() *exec.Cmd {
	appArgs := []string{"get", "-t"}
	if verbose {
		appArgs = append(appArgs, "-v")
	}
	appArgs = append(appArgs, "./...")
	return exec.Command("go", appArgs...)
}

const nMain = `package main

import (
	"log"
	"net/http"

	"{{.actionsPath}}"
)

func main() {
	log.Fatal(http.ListenAndServe(":3000", actions.App()))
}

`
const nApp = `package actions

import (
	"net/http"

	"github.com/markbates/buffalo"
	{{if .withPop -}}
	"github.com/markbates/buffalo/middleware"
	"{{.modelsPath}}"
	{{end -}}
)

func App() http.Handler {
	a := buffalo.Automatic(buffalo.Options{
		Env: "development",
	})

	{{if .withPop -}}
	a.Use(middleware.PopTransaction(models.DB))
	{{end -}}

	a.ServeFiles("/assets", assetsPath())
	a.GET("/", HomeHandler)

	return a
}
`

const nRender = `package actions

import (
	"net/http"
	"path"
	"runtime"

	"github.com/markbates/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
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

const nHomeHandler = `package actions

import "github.com/markbates/buffalo"

func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}
`

const nHomeHandlerTest = `package actions_test

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

const nIndexHTML = `<h1>Welcome to Buffalo!</h1>`

const nApplicationHTML = `<html>
<head>
  <meta charset="utf-8">
  <title>Buffalo - {{ .name }}</title>
  <link rel="stylesheet" href="/assets/application.css" type="text/css" media="all" />
</head>
<body>
  {{"{{"}} yield {{"}}"}}

	{{if .withJQuery -}}
  <script src="/assets/jquery.js" type="text/javascript" charset="utf-8"></script>
	{{end -}}
  <script src="/assets/application.js" type="text/javascript" charset="utf-8"></script>
</body>
</html>
`

const nApplicationCSS = `body {
  font-family: helvetica;
}
`

const nGitignore = `vendor/
**/*.log
**/*.sqlite
bin/
node_modules/
{{ .name }}
`

const nGriftRoutes = `package grifts

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

const nRefresh = `app_root: .
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
enable_colors: true
log_name: buffalo
`
