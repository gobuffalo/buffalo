package cmd

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/destroy"
	"github.com/gobuffalo/pop/soda/cmd"
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
	decorate("destroy", destroyCmd)
	c.AddCommand(destroyCmd)

	c.Use = "db"
	decorate("db", c)
	RootCmd.AddCommand(c)
}
