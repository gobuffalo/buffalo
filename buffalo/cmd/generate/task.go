package generate

import (
	"context"
	"os"

	"github.com/gobuffalo/buffalo/genny/grift"
	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/gogen"
	"github.com/spf13/cobra"
)

var taskOptions = struct {
	dryRun bool
	*grift.Options
}{
	Options: &grift.Options{},
}

//TaskCmd is the command called with the generate grift cli.
var TaskCmd = &cobra.Command{
	Use:     "task [name]",
	Aliases: []string{"t", "grift"},
	Short:   "Generate a grift task",
	RunE: func(cmd *cobra.Command, args []string) error {
		run := genny.WetRunner(context.Background())
		if taskOptions.dryRun {
			run = genny.DryRunner(context.Background())
		}

		opts := taskOptions.Options
		opts.Args = args
		g, err := grift.New(opts)
		if err != nil {
			return err
		}
		run.With(g)

		pwd, _ := os.Getwd()
		g, err = gogen.Fmt(pwd)
		if err != nil {
			return err
		}
		run.With(g)

		return run.Run()
	},
}

func init() {
	TaskCmd.Flags().BoolVarP(&taskOptions.dryRun, "dry-run", "d", false, "dry run of the generator")
}
