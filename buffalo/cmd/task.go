package cmd

import (
	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

// task command is a forward to grift tasks
var taskCommand = &cobra.Command{
	Use:                "task",
	Aliases:            []string{"t", "tasks"},
	Short:              "Runs your grift tasks",
	DisableFlagParsing: true,
	RunE: func(c *cobra.Command, args []string) error {
		return grifts.Run("buffalo task", args)
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	decorate("task", taskCommand)
	RootCmd.AddCommand(taskCommand)
}
