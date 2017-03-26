package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of the buffalo binary
const Version = "0.8.0.dev"

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of buffalo",
	Long:  `All software has versions.  This is buffalo's.`,
	Run: func(c *cobra.Command, args []string) {
		fmt.Println("Buffalo version %s", Version)
	},
}
