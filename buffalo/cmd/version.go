package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Version is the current version of the buffalo binary
const Version = "v0.11.1"

func init() {
	decorate("version", versionCmd)
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of buffalo",
	Long:  `All software has versions.  This is buffalo's.`,
	Run: func(c *cobra.Command, args []string) {
		logrus.Infof("Buffalo version is: %s\n", Version)
	},
	// needed to override the root level pre-run func
	PersistentPreRunE: func(c *cobra.Command, args []string) error {
		return nil
	},
}
