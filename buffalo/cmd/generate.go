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
	generate.Version = Version
	generateCmd.AddCommand(generate.ResourceCmd)
	generateCmd.AddCommand(generate.ActionCmd)
	generateCmd.AddCommand(generate.DockerCmd)
	generateCmd.AddCommand(generate.TaskCmd)
	generateCmd.AddCommand(generate.MailCmd)
	decorate("generate", generateCmd)

	RootCmd.AddCommand(generateCmd)
}
