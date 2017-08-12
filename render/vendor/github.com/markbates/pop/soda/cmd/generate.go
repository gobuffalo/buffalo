package cmd

import (
	"github.com/markbates/pop/soda/cmd/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"g"},
}

func init() {
	generateCmd.AddCommand(generate.ConfigCmd)
	generateCmd.AddCommand(generate.FizzCmd)
	generateCmd.AddCommand(generate.SQLCmd)
	generateCmd.AddCommand(generate.ModelCmd)
	RootCmd.AddCommand(generateCmd)
}
