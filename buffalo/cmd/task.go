package cmd

import (
	"errors"
	"os"

	grifts "github.com/markbates/grift/cmd"
	"github.com/spf13/cobra"
)

// task command is a forward to grift tasks
var taskCommand = &cobra.Command{
	Use:     "task",
	Aliases: []string{"t", "tasks"},
	Short:   "Runs your grift tasks",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := os.Stat("grifts")
		if err != nil {
			return errors.New("seems there is no grifts folder on your current directory, please ensure you're inside your buffalo app root")
		}

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
