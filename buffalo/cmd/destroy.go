package cmd

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/destroy"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "Allows to destroy generated code.",
	Aliases: []string{"d"},
}

func init() {
	destroyCmd.AddCommand(destroy.ResourceCmd)
	destroyCmd.AddCommand(destroy.ActionCmd)

	destroy.ResourceCmd.Flags().BoolVarP(&destroy.YesToAll, "yes", "y", false, "confirms all beforehand")
	destroy.ModelCmd.Flags().BoolVarP(&destroy.YesToAll, "yes", "y", false, "confirms all beforehand")
	destroy.ActionCmd.Flags().BoolVarP(&destroy.YesToAll, "yes", "y", false, "confirms all beforehand")

	decorate("destroy", destroyCmd)
	RootCmd.AddCommand(destroyCmd)
}
