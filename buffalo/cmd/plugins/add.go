package plugins

import (
	"context"
	"os"
	"path"
	"strings"

	"github.com/gobuffalo/buffalo/genny/add"
	"github.com/gobuffalo/buffalo/internal/takeon/github.com/markbates/errx"
	"github.com/gobuffalo/buffalo/plugins/plugdeps"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/meta"
	"github.com/spf13/cobra"
)

var addOptions = struct {
	dryRun    bool
	buildTags []string
}{}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "adds plugins to config/buffalo-plugins.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		run := genny.WetRunner(context.Background())
		if addOptions.dryRun {
			run = genny.DryRunner(context.Background())
		}

		app := meta.New(".")
		plugs, err := plugdeps.List(app)
		if err != nil && (errx.Unwrap(err) != plugdeps.ErrMissingConfig) {
			return err
		}

		tags := app.BuildTags("", addOptions.buildTags...)
		for _, a := range args {
			a = strings.TrimSpace(a)
			bin := path.Base(a)
			plug := plugdeps.Plugin{
				Binary: bin,
				GoGet:  a,
				Tags:   tags,
			}
			if _, err := os.Stat(a); err == nil {
				plug.Local = a
				plug.GoGet = ""
			}
			plugs.Add(plug)
		}
		g, err := add.New(&add.Options{
			App:     app,
			Plugins: plugs.List(),
		})
		if err != nil {
			return err
		}
		run.With(g)

		return run.Run()
	},
}

func init() {
	addCmd.Flags().BoolVarP(&addOptions.dryRun, "dry-run", "d", false, "dry run")
	addCmd.Flags().StringSliceVarP(&addOptions.buildTags, "tags", "t", []string{}, "build tags")
}
