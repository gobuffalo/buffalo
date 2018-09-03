package cmd

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/destroy"
	"github.com/spf13/cobra"
)

// DestroyCmd destroys generated code
var DestroyCmd = &cobra.Command{
	Use:     "destroy",
	Short:   "destroys generated code",
	Aliases: []string{"d"},
}

func init() {
	DestroyCmd.AddCommand(destroy.ResourceCmd)
	DestroyCmd.AddCommand(destroy.ActionCmd)
	DestroyCmd.AddCommand(destroy.MailerCmd)

	DestroyCmd.PersistentFlags().BoolVarP(&destroy.YesToAll, "yes", "y", false, "confirms all beforehand")

	decorate("destroy", DestroyCmd)
	RootCmd.AddCommand(DestroyCmd)
}
