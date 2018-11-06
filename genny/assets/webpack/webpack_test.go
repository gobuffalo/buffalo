package webpack

import (
	"strconv"
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
	r.Len(res.Files, 9)

	c := res.Commands[0]
	r.Equal("npm install --no-progress --save", strings.Join(c.Args, " "))

	f := res.Files[0]
	r.Equal(".babelrc", f.Name())

	f = res.Files[1]
	r.Equal("assets/css/application.scss", f.Name())

	f = res.Files[2]
	r.Equal("assets/images/favicon.ico", f.Name())

	f = res.Files[3]
	r.Equal("assets/images/logo.svg", f.Name())

	f = res.Files[4]
	r.Equal("assets/js/application.js", f.Name())

	f = res.Files[5]
	r.Equal("package.json", f.Name())
	r.Contains(f.String(), `"bootstrap": "4.1.1",`)

	f = res.Files[6]
	r.Equal("public/assets/.keep", f.Name())

	f = res.Files[7]
	r.Equal("templates/application.html", f.Name())

	f = res.Files[8]
	r.Equal("webpack.config.js", f.Name())
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
	r.Len(res.Files, 9)

	c := res.Commands[0]
	r.Equal("yarnpkg install --no-progress --save", strings.Join(c.Args, " "))
}

func Test_Webpack_Updates_Layout(t *testing.T) {
	table := []struct {
		v   int
		css string
	}{
		{3, bs3},
		{4, bs4},
	}

	for _, tt := range table {
		t.Run(strconv.Itoa(tt.v), func(st *testing.T) {
			r := require.New(st)
			run := runner()

			run.WithNew(New(&Options{
				Bootstrap: tt.v,
			}))

			r.NoError(run.Run())

			res := run.Results()

			f, err := res.Find("templates/application.html")
			r.NoError(err)

			body := f.String()
			r.Contains(body, "</title>\n"+tt.css)
			r.Contains(body, `<%= stylesheetTag("application.css") %>`)
		})
	}
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
