package webpack

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func runner() *genny.Runner {
	run := gentest.NewRunner()
	run.Disk.Add(genny.NewFileS("templates/application.html", layout))
	run.LookPathFn = func(s string) (string, error) {
		return s, nil
	}
	return run
}

func Test_Webpack_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := runner()

	run.With(g)
	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 1)
	c := res.Commands[0]
	r.Equal("npm install --no-progress --save", strings.Join(c.Args, " "))

	files := []string{
		".babelrc",
		"assets/css/_buffalo.scss",
		"assets/css/application.scss",
		"assets/images/favicon.ico",
		"assets/images/logo.svg",
		"assets/js/application.js",
		"package.json",
		"public/assets/.keep",
		"templates/application.html",
		"webpack.config.js",
	}
	r.Len(res.Files, len(files))
	for i, f := range res.Files {
		r.Equal(files[i], f.Name())
	}

	f, err := res.Find("package.json")
	r.NoError(err)
	r.Contains(f.String(), `"bootstrap": "4.`)

}

func Test_Webpack_New_WithYarn(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		App: meta.App{WithYarn: true},
	})
	r.NoError(err)

	run := runner()
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 10)

	c := res.Commands[0]
	r.Equal("yarnpkg install --no-progress --save", strings.Join(c.Args, " "))
}

const layout = `<!DOCTYPE html>
<html>
  <head>
    <title>Buffalo - Foo</title>
    <%= stylesheetTag("buffalo.css") %>
    <%= stylesheetTag("application.css") %>
  </head>
  <body>
  </body>
</html>
`
