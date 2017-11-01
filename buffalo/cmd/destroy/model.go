package destroy

import (
	"errors"

	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

//ModelCmd destroys a passed model
var ModelCmd = &cobra.Command{
	Use: "model [name]",
	//Example: "resource cars",
	Aliases: []string{"m"},
	Short:   "Destroys model files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you need to provide a valid model name in order to destroy it")
		}

		name := args[0]
		fileName := inflect.Pluralize(inflect.Underscore(name))

		removeModel(name)
		removeMigrations(fileName)

		return nil
	},
}
