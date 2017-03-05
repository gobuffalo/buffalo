package generate

import (
	"errors"
	"os"

	"github.com/gobuffalo/buffalo/generators/action"
	"github.com/markbates/gentronics"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

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

		data := gentronics.Data{
			"filename":  inflect.Underscore(name),
			"namespace": inflect.Camelize(name),
		}

		g, err := action.New(name, args[1:], data)
		if err != nil {
			return err
		}

		return g.Run(".", data)
	},
}
