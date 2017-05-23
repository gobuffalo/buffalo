package generate

import (
	"errors"
	"os"

	"github.com/gobuffalo/buffalo/generators/action"
	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

//SkipActionTemplate indicates whether we generator should not generate the view layer when generating actions.
var SkipActionTemplate = false

//ActionMethod is the method generated action will be binded to.
var ActionMethod = "GET"

//ActionCmd is the cmd that generates actions.
var ActionCmd = &cobra.Command{
	Use:     "action [name] [actionName...]",
	Aliases: []string{"a", "actions"},
	Short:   "Generates new action(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return errors.New("you should provide action name and handler name at least")
		}

		if _, err := os.Stat("actions"); err != nil {
			return errors.New("actions directory not found, ensure you're inside your buffalo folder")
		}

		name := args[0]

		data := makr.Data{
			"filename":     inflect.Underscore(name),
			"namespace":    inflect.Camelize(name),
			"method":       ActionMethod,
			"skipTemplate": SkipActionTemplate,
		}

		g, err := action.New(name, args[1:], data)
		if err != nil {
			return err
		}

		return g.Run(".", data)
	},
}
