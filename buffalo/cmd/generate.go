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
	generateCmd.AddCommand(generate.WebpackCmd)
	generateCmd.AddCommand(generate.ActionCmd)
	generateCmd.AddCommand(generate.DockerCmd)
	generateCmd.AddCommand(generate.TaskCmd)
	decorate("generate", generateCmd)

	generate.ResourceCmd.Flags().BoolVarP(&generate.SkipResourceMigration, "skip-migration", "s", false, "sets resource generator not-to add model migration")
	generate.ResourceCmd.Flags().BoolVarP(&generate.SkipResourceModel, "skip-model", "", false, "makes resource generator not to generate model nor migrations")
	generate.ActionCmd.Flags().BoolVarP(&generate.SkipActionTemplate, "skip-template", "", false, "makes resource generator not to generate template for actions")

	generate.ActionCmd.Flags().StringVarP(&generate.ActionMethod, "method", "m", "GET", "allows to set a different method for the action being generated.")
	generate.ResourceCmd.Flags().StringVarP(&generate.UseResourceModel, "use-model", "u", "", "generates crud options for a model")
	generate.ResourceCmd.Flags().StringVarP(&generate.ModelName, "model-name", "n", "", "allows to define a different model name for the resource being generated.")
	generate.ResourceCmd.Flags().StringVarP(&generate.ResourceMimeType, "type", "", "html", "sets the resource type (html or json)")

	RootCmd.AddCommand(generateCmd)
}
