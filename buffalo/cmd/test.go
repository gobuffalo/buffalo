package cmd

import (
	"os"

	tt "github.com/markbates/tt/cmd"
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:                "test",
	Short:              "Runs the tests for your Buffalo app",
	DisableFlagParsing: true,
	Run: func(c *cobra.Command, args []string) {
		os.Setenv("GO_ENV", "test")
		tt.Run(tt.GoBuilder(args))
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
}
