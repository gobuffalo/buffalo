package generate

import (
	"context"
	"fmt"

	"github.com/gobuffalo/buffalo/genny/actions"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/logger"
	"github.com/spf13/cobra"
)

var actionOptions = struct {
	*actions.Options
	dryRun  bool
	verbose bool
}{
	Options: &actions.Options{},
}

//ActionCmd is the cmd that generates actions.
var ActionCmd = &cobra.Command{
	Use:     "action [name] [handler name...]",
	Aliases: []string{"a", "actions"},
	Short:   "Generate new action(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you must provide a name")
		}
		actionOptions.Name = args[0]
		if len(args) == 1 {
			return fmt.Errorf("you must provide at least one action name")
		}
		actionOptions.Actions = args[1:]

		ctx := context.Background()
		run := genny.WetRunner(ctx)

		if actionOptions.dryRun {
			run = genny.DryRunner(ctx)
		}

		if actionOptions.verbose {
			run.Logger = logger.New(logger.DebugLevel)
		}

		opts := actionOptions.Options
		run.WithNew(actions.New(opts))
		return run.Run()
	},
}

func init() {
	ActionCmd.Flags().BoolVarP(&actionOptions.SkipTemplates, "skip-template", "", false, "skip generation of templates for action(s)")
	ActionCmd.Flags().BoolVarP(&actionOptions.dryRun, "dry-run", "d", false, "dry run")
	ActionCmd.Flags().BoolVarP(&actionOptions.verbose, "verbose", "v", false, "verbosely run the generator")
	ActionCmd.Flags().StringVarP(&actionOptions.Method, "method", "m", "GET", "change the HTTP method for the generate action(s)")
}
