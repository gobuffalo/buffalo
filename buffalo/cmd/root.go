package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gobuffalo/events"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var anywhereCommands = []string{"new", "version", "info", "help"}

// RootCmd is the hook for all of the other commands in the buffalo binary.
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "buffalo",
	Short:         "Helps you build your Buffalo applications that much easier!",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := events.LoadPlugins(); err != nil {
			return err
		}
		isFreeCommand := false
		for _, freeCmd := range anywhereCommands {
			if freeCmd == cmd.Name() {
				isFreeCommand = true
			}
		}

		if isFreeCommand {
			return nil
		}

		if !insideBuffaloProject() {
			return errors.New("you need to be inside your buffalo project path to run this command")
		}

		if err := startDockerCompose(); err != nil {
			return errors.WithStack(err)
		}

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Errorf("Error: %s\n\n", err)
		os.Exit(-1)
	}
}

func init() {
	decorate("root", RootCmd)
}

func insideBuffaloProject() bool {
	if _, err := os.Stat(".buffalo.dev.yml"); err != nil {
		return false
	}

	return true
}

func startDockerCompose() error {
	if _, err := os.Stat("docker-compose.yml"); err != nil {
		return nil
	}

	if _, err := exec.LookPath("docker-compose"); err != nil {
		if err != nil {
			return errors.New("This application require docker-compose and we could not find it installed on your system")
		}
	}

	fmt.Println("Start docker-compose")

	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
