package cmd

import (
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	decorate("version", versionCmd)
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Run: func(c *cobra.Command, args []string) {
		logrus.Infof("Buffalo version is: %s\n", runtime.Version)
	},
	// needed to override the root level pre-run func
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
