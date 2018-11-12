package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/gobuffalo/events"
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

		return nil
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		if strings.Contains(err.Error(), dbNotFound) || strings.Contains(err.Error(), popNotFound) {
			logrus.Errorf(popInstallInstructions)
			os.Exit(-1)
		}
		logrus.Errorf("Error: %s\n\n", err)
		os.Exit(-1)
	}
}

const dbNotFound = `unknown command "db"`
const popNotFound = `unknown command "pop"`
const popInstallInstructions = `Pop support has been moved to the https://github.com/gobuffalo/buffalo-pop plugin.

Go Get Installation:

	$ go get github.com/gobuffalo/buffalo-pop

Buffalo Plugins Installation*:

	$ buffalo plugins install github.com/gobuffalo/buffalo-pop

* Requires https://github.com/gobuffalo/buffalo-plugins installed.
`

func init() {
	decorate("root", RootCmd)
}

func insideBuffaloProject() bool {
	if _, err := os.Stat(".buffalo.dev.yml"); err != nil {
		return false
	}

	return true
}
