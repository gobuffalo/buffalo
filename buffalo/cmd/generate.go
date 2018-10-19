package cmd

import (
	"github.com/gobuffalo/buffalo/buffalo/cmd/generate"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate application components",
	Aliases: []string{"g"},
}

func init() {
	generateCmd.AddCommand(generate.ResourceCmd)
	generateCmd.AddCommand(generate.ActionCmd)
	generateCmd.AddCommand(generate.TaskCmd)
	generateCmd.AddCommand(generate.MailCmd)
	decorate("generate", generateCmd)

	RootCmd.AddCommand(generateCmd)
}
