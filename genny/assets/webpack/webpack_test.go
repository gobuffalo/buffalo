package webpack

import (
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_Webpack_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 8)

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
	r.Equal("webpack.config.js", f.Name())
}

func Test_Webpack_New_WithYarn(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		App: meta.App{WithYarn: true},
	})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 1)
	r.Len(res.Files, 8)

	c := res.Commands[0]
	r.Equal("yarnpkg install --no-progress --save", strings.Join(c.Args, " "))
}
