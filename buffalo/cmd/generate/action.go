package generate

import (
	"github.com/pkg/errors"

	"github.com/gobuffalo/buffalo/generators/action"
	"github.com/gobuffalo/makr"
	"github.com/spf13/cobra"
)

var actionOptions = struct {
	SkipTemplate bool
	Method       string
}{}

//ActionCmd is the cmd that generates actions.
var ActionCmd = &cobra.Command{
	Use:     "action [name] [actionName...]",
	Aliases: []string{"a", "actions"},
	Short:   "Generates new action(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := action.New(args...)
		if err != nil {
			return errors.WithStack(err)
		}
		a.SkipTemplate = actionOptions.SkipTemplate
		a.Method = actionOptions.Method

		data := makr.Data{}

		return a.Run(".", data)
	},
}

func init() {
	ActionCmd.Flags().BoolVarP(&actionOptions.SkipTemplate, "skip-template", "", false, "skip generation of templates for action(s)")
	ActionCmd.Flags().StringVarP(&actionOptions.Method, "method", "m", "GET", "change the HTTP method for the generate action(s)")
}
