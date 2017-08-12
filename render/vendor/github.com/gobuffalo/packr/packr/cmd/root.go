package cmd

import (
	"context"
	"os"

	"github.com/gobuffalo/packr/builder"
	"github.com/spf13/cobra"
)

var input string

var rootCmd = &cobra.Command{
	Use:   "packr",
	Short: "compiles static files into Go files",
	RunE: func(cmd *cobra.Command, args []string) error {
		b := builder.New(context.Background(), input)
		return b.Run()
	},
}

func init() {
	rootCmd.Flags().StringVarP(&input, "input", "i", ".", "path to scan for packr Boxes")
}

// Execute the commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
