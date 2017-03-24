package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo/buffalo/cmd/generate"
	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd is the hook for all of the other commands in the buffalo binary.
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "buffalo",
	Short:         "Helps you build your Buffalo applications that much easier!",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Buffalo version %s\n\n", Version)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Error: %s\n\n", err)
		os.Exit(-1)
	}
}

func goInstall(pkg string) *exec.Cmd {
	return generate.GoInstall(pkg)
}

func goGet(pkg string) *exec.Cmd {
	return generate.GoGet(pkg)
}
