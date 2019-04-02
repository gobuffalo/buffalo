package standard

import (
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
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

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(nil)
	r.NoError(err)

	run := runner()
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)

	files := []string{
		"public/assets/application.css",
		"public/assets/application.js",
		"public/assets/buffalo.css",
		"public/assets/images/favicon.ico",
	}
	r.Len(res.Files, len(files))
	for i, f := range res.Files {
		r.Equal(files[i], f.Name())
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
