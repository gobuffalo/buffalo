package cmd

import (
	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

func init() {
	decorate("routes", routesCmd)
	RootCmd.AddCommand(routesCmd)
}

var routesCmd = &cobra.Command{
	Use:   "routes",
	Short: "Print out all defined routes",
	RunE: func(c *cobra.Command, args []string) error {
		return grifts.Run("buffalo task", []string{"routes"})
	},
}
