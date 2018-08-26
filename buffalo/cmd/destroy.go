package cmd

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/destroy"
	"github.com/spf13/cobra"
)

var DestroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "Allows to destroy generated code.",
	Aliases: []string{"d"},
}

func init() {
	DestroyCmd.AddCommand(destroy.ResourceCmd)
	DestroyCmd.AddCommand(destroy.ActionCmd)

	DestroyCmd.PersistentFlags().BoolVarP(&destroy.YesToAll, "yes", "y", false, "confirms all beforehand")

	decorate("destroy", DestroyCmd)
	RootCmd.AddCommand(DestroyCmd)
}
