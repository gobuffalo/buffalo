package cmd

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "A collection of generators to make life easier",
	Aliases: []string{"g"},
}

func init() {
	generateCmd.AddCommand(generate.ResourceCmd)
	generateCmd.AddCommand(generate.GothCmd)
	generateCmd.AddCommand(generate.WebpackCmd)
	generateCmd.AddCommand(generate.ActionCmd)
	RootCmd.AddCommand(generateCmd)
}
