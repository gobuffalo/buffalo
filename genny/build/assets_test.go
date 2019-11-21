package build

import (
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/internal/envx"
	"github.com/stretchr/testify/require"
)

func Test_assets(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		WithAssets: true,
	}
	r.NoError(opts.Validate())
	opts.App.WithNodeJs = true
	opts.App.PackageJSON.Scripts = map[string]string{
		"build": "webpack -p --progress",
	}

	run := cokeRunner()
	run.WithNew(assets(opts))

	os.Setenv("NODE_ENV", "")
	ne := envx.Get("NODE_ENV", "")
	r.Empty(ne)
	r.NoError(run.Run())

	ne = envx.Get("NODE_ENV", "")
	r.NotEmpty(ne)
	r.Equal(opts.Environment, ne)

	res := run.Results()

	cmds := []string{"npm run build"}
	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}
}

func Test_assets_Archived(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		WithAssets:    true,
		ExtractAssets: true,
	}
	r.NoError(opts.Validate())

	run := cokeRunner()
	opts.Root = run.Root
	run.WithNew(assets(opts))
	r.NoError(run.Run())

	res := run.Results()

	cmds := []string{}
	r.Len(res.Commands, len(cmds))
	for i, c := range res.Commands {
		r.Equal(cmds[i], strings.Join(c.Args, " "))
	}

	// r.Len(res.Files, 1)

	f, err := res.Find("actions/app.go")
	r.NoError(err)
	r.Contains(f.String(), `// app.ServeFiles("/"`)
}
