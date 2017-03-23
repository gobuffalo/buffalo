package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// task command is a forward to grift tasks
var taskCommand = &cobra.Command{
	Use:     "task",
	Aliases: []string{"t", "tasks"},
	Short:   "Runs your grift tasks",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := exec.LookPath("grift")
		if err != nil {
			return errors.New("we could not find \"grift\" in your path.\n You must first install \"grift\" in order to use the Buffalo console:\n\n $ go get github.com/markbates/grift")
		}

		_, err = os.Stat("grifts")
		if err != nil {
			return errors.New("seems there is no grift folder on your current directory, please ensure you're inside your buffalo app root")
		}

		cmd := exec.Command("grift", args...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()

	},
}

func init() {
	RootCmd.AddCommand(taskCommand)
}
