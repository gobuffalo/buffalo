package standard

import (
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := gentest.NewRunner()
	run.Disk.Add(genny.NewFileS("templates/application.plush.html", layout))
	run.LookPathFn = func(s string) (string, error) {
		return s, nil
	}

	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)

	files := []string{
		"public/assets/application.css",
		"public/assets/application.js",
		"public/assets/buffalo.css",
		"public/assets/images/favicon.ico",
		"public/assets/images/logo.svg",
		"templates/application.plush.html",
	}

	r.Len(res.Files, len(files))
	for i, f := range res.Files {
		r.Equal(files[i], f.Name())
	}

	layout, ferr := res.Find("templates/application.plush.html")
	r.NoError(ferr)

	r.Contains(layout.String(), "href=\"https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css\"")
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
