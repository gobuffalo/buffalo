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

	generate.ResourceCmd.Flags().BoolVarP(&generate.SkipResourceMigration, "skip-migration", "s", false, "sets resource generator not-to add model migration")
	generate.ResourceCmd.Flags().BoolVarP(&generate.SkipResourceModel, "skip-model", "", false, "makes resource generator not to generate model nor migrations")
	generate.ResourceCmd.Flags().BoolVarP(&generate.UseResourceModel, "use-model", "u", false, "generates crud options for a model")

	RootCmd.AddCommand(generateCmd)
}
