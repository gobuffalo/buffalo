package plugins

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo/genny/plugins/plugin"
	"github.com/gobuffalo/buffalo/genny/plugins/plugin/with"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen"
	"github.com/gobuffalo/licenser/genny/licenser"
	"github.com/gobuffalo/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "plugin",
	Short: "generates a new buffalo plugin",
	Long:  "buffalo generate plugin github.com/foo/buffalo-bar",
	RunE: func(cmd *cobra.Command, args []string) error {
		popts := &plugin.Options{
			Author:    viper.GetString("author"),
			ShortName: viper.GetString("short-name"),
			License:   viper.GetString("license"),
		}
		if len(args) > 0 {
			popts.PluginPkg = args[0]
		}

		r := genny.WetRunner(context.Background())
		if viper.GetBool("dry-run") {
			r = genny.DryRunner(context.Background())
		}

		popts.Root = filepath.Join(envy.GoPath(), "src")

		gg, err := plugin.New(popts)
		if err != nil {
			return err
		}
		r.Root = popts.Root
		r.WithRun(genny.Force(r.Root, viper.GetBool("force")))
		r.WithGroup(gg)

		if viper.GetBool("with-gen") {
			gg, err := with.GenerateCmd(popts)
			if err != nil {
				return err
			}
			r.WithGroup(gg)
		}

		g, err := gogen.Fmt(r.Root)
		if err != nil {
			return err
		}
		r.With(g)

		if viper.GetBool("verbose") {
			r.Logger = logger.New(logger.DebugLevel)
		}
		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolP("dry-run", "d", false, "run the generator without creating files or running commands")
	generateCmd.Flags().BoolP("verbose", "v", false, "turn on verbose logging")
	generateCmd.Flags().Bool("with-gen", false, "creates a generator plugin")
	generateCmd.Flags().BoolP("force", "f", false, "will delete the target directory if it exists")
	generateCmd.Flags().StringP("author", "a", "", "author's name")
	generateCmd.Flags().StringP("license", "l", "mit", fmt.Sprintf("choose a license from: [%s]", strings.Join(licenser.Available, ", ")))
	generateCmd.Flags().StringP("short-name", "s", "", "a 'short' name for the package")
	viper.BindPFlags(generateCmd.Flags())
}
