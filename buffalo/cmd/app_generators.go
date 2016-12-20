package cmd

import (
	"os/exec"

	"github.com/markbates/buffalo/buffalo/cmd/generate"
	"github.com/markbates/gentronics"
)

func newAppGenerator(data gentronics.Data) *gentronics.Generator {
	g := gentronics.New()
	g.Add(gentronics.NewFile("main.go", nMain))
	g.Add(gentronics.NewFile("Procfile", nProcfile))
	g.Add(gentronics.NewFile("Procfile.development", nProcfileDev))
	g.Add(gentronics.NewFile(".buffalo.dev.yml", nRefresh))
	g.Add(gentronics.NewFile("actions/app.go", nApp))
	g.Add(gentronics.NewFile("actions/home.go", nHomeHandler))
	g.Add(gentronics.NewFile("actions/home_test.go", nHomeHandlerTest))
	g.Add(gentronics.NewFile("actions/render.go", nRender))
	g.Add(gentronics.NewFile("grifts/routes.go", nGriftRoutes))
	g.Add(gentronics.NewFile("templates/index.html", nIndexHTML))
	g.Add(gentronics.NewFile("templates/application.html", nApplicationHTML))
	if skipWebpack {
		g.Add(gentronics.NewFile("assets/js/application.js", ""))
		g.Add(gentronics.NewFile("assets/css/application.css", ""))
	}
	g.Add(gentronics.NewFile(".gitignore", nGitignore))
	g.Add(gentronics.NewCommand(goGet("github.com/markbates/refresh/...")))
	g.Add(gentronics.NewCommand(goInstall("github.com/markbates/refresh")))
	g.Add(gentronics.NewCommand(goGet("github.com/markbates/grift/...")))
	g.Add(gentronics.NewCommand(goInstall("github.com/markbates/grift")))
	g.Add(gentronics.NewCommand(goGet("github.com/motemen/gore")))
	g.Add(gentronics.NewCommand(goInstall("github.com/motemen/gore")))
	g.Add(generate.NewWebpackGenerator(data))
	g.Add(newSodaGenerator())
	g.Add(gentronics.NewCommand(appGoGet()))
	g.Add(generate.Fmt)
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
	"fmt"
	"log"
	"net/http"
	"os"

	"{{actionsPath}}"
	"github.com/markbates/going/defaults"
)

func main() {
	port := defaults.String(os.Getenv("PORT"), "3000")
	log.Printf("Starting {{name}} on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), actions.App()))
}

`
const nApp = `package actions

import (
	"os"

	"github.com/markbates/buffalo"
	{{#if withPop }}
	"github.com/markbates/buffalo/middleware"
	"{{modelsPath}}"
	{{/if}}
	"github.com/markbates/going/defaults"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = defaults.String(os.Getenv("GO_ENV"), "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env: ENV,
		})

		{{#if withPop }}
		app.Use(middleware.PopTransaction(models.DB))
		{{/if}}

		app.ServeFiles("/assets", assetsPath())
		app.GET("/", HomeHandler)
	}

	return app
}
`

const nRender = `package actions

import (
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/markbates/buffalo/render"
	"github.com/markbates/buffalo/render/resolvers"
)

var r *render.Engine
var resolver = &resolvers.RiceBox{
	Box: rice.MustFindBox("../templates"),
}

func init() {
	r = render.New(render.Options{
		HTMLLayout:     "application.html",
		CacheTemplates: ENV == "production",
		FileResolver:   resolver,
	})
}

func assetsPath() http.FileSystem {
	box := rice.MustFindBox("../assets")
	return box.HTTPBox()
}
`

const nHomeHandler = `package actions

import "github.com/markbates/buffalo"

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}
`

const nHomeHandlerTest = `package actions_test

import (
	"testing"

	"{{actionsPath}}"
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
  <title>Buffalo - {{ titleName }}</title>
	{{#if withWebpack}}
		<link rel="stylesheet" href="/assets/dist/application.css" type="text/css" media="all" />
	{{else}}
		<link rel="stylesheet" href="/assets/css/application.css" type="text/css" media="all" />
	{{/if}}
</head>
<body>

	\{{ yield }}
	{{#if withWebpack}}
		<script src="/assets/dist/application.js" type="text/javascript" charset="utf-8"></script>
	{{else}}
		<script src="/assets/js/application.js" type="text/javascript" charset="utf-8"></script>
	{{/if}}
</body>
</html>
`

const nGitignore = `vendor/
**/*.log
**/*.sqlite
bin/
node_modules/
.sass-cache/
assets/dist/
{{ name }}
`

const nGriftRoutes = `package grifts

import (
	"os"

	. "github.com/markbates/grift/grift"
	"{{actionsPath}}"
	"github.com/olekukonko/tablewriter"
)

var _ = Add("routes", func(c *Context) error {
	a := actions.App()
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
- logs
- assets
- grifts
- tmp
- node_modules
- .sass-cache
included_extensions:
- .go
- .html
- .md
build_path: /tmp
build_delay: 200ns
binary_name: {{name}}-build
command_flags: []
enable_colors: true
log_name: buffalo
`

const nProcfile = `web: {{name}}`
const nProcfileDev = `web: buffalo dev
{{#if withWebpack}}
assets: webpack --watch
{{/if}}
`
