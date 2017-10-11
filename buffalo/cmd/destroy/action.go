package destroy

import (
	"errors"

	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

//ActionCmd destroys passed action file
var ActionCmd = &cobra.Command{
	Use: "action [name]",
	//Example: "resource cars",
	Aliases: []string{"a"},
	Short:   "Destroys action files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("you need to provide a valid action file name in order to destroy it")
		}

		name := args[0]

		//Generated actions keep the same name (not plural).
		fileName := inflect.Underscore(name)

		removeActions(fileName)
		return nil
	},
}
