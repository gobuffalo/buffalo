package cmd

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/destroy"
	"github.com/markbates/pop/soda/cmd"
	"github.com/spf13/cobra"
)

func init() {

	var destroyCmd = &cobra.Command{
		Use:     "destroy",
		Short:   "Allows to destroy generated code.",
		Aliases: []string{"d"},
	}

	c := cmd.RootCmd
	destroyCmd.AddCommand(destroy.ModelCmd)
	c.AddCommand(destroyCmd)

	c.Use = "db"
	RootCmd.AddCommand(c)
}
