package with

import (
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/genny/plugin"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/gogen/gomods"
	"github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

var gBox = packr.New("generate_test", "../../plugin/templates")

func Test_GenerateCmd(t *testing.T) {
	r := require.New(t)

	err := gomods.Disable(func() error {
		opts := &plugin.Options{
			PluginPkg: "github.com/foo/buffalo-bar",
			Year:      1999,
			Author:    "Homer Simpson",
			ShortName: "bar",
		}

		run := gentest.NewRunner()
		run.Disk.Add(genny.NewFile("cmd/available.go", strings.NewReader(availgo)))

		gg, err := GenerateCmd(opts)
		r.NoError(err)
		run.WithGroup(gg)

		r.NoError(run.Run())

		res := run.Results()
		r.Len(res.Commands, 0)
		r.Len(res.Files, 7)

		f := res.Files[0]
		r.Equal("cmd/available.go", f.Name())
		r.Contains(f.String(), `Available.Add("generate", generateCmd)`)

		f = res.Files[1]
		r.Equal("cmd/generate.go", f.Name())
		r.Contains(f.String(), opts.PluginPkg+"/genny/")

		f = res.Files[2]
		r.Equal("genny/bar/bar.go", f.Name())
		r.Contains(f.String(), "package bar")
		r.Contains(f.String(), "func New(opts *Options) (*genny.Generator, error)")

		f = res.Files[3]
		r.Equal("genny/bar/bar_test.go", f.Name())

		f = res.Files[4]
		r.Equal("genny/bar/options.go", f.Name())
		r.Contains(f.String(), "package bar")
		r.Contains(f.String(), "type Options struct {")

		f = res.Files[5]
		r.Equal("genny/bar/options_test.go", f.Name())

		f = res.Files[6]
		r.Equal("genny/bar/templates/example.txt", f.Name())
		return nil
	})
	r.NoError(err)
}

const availgo = `package cmd

import (
	"github.com/gobuffalo/buffalo/plugins/plugcmds"
	"github.com/spf13/cobra"
)

var Available = plugcmds.NewAvailable()

var pluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "tools for working with buffalo plugins",
}

func init() {
	Available.Add("root", pluginsCmd)
	Available.Mount(rootCmd)
}`
