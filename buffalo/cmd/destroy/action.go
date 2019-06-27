package destroy

import (
	"fmt"

	"github.com/gobuffalo/flect"
	"github.com/spf13/cobra"
)

//ActionCmd destroys passed action file
var ActionCmd = &cobra.Command{
	Use: "action [name]",
	//Example: "resource cars",
	Aliases: []string{"a"},
	Short:   "Destroy action files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("you need to provide a valid action file name in order to destroy it")
		}

		name := args[0]

		//Generated actions keep the same name (not plural).
		fileName := flect.Underscore(name)

		removeActions(fileName)
		return nil
	},
}
