package cmd

import (
	"github.com/markbates/pop/soda/cmd/schema"
	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Tools for working with your database schema",
}

func init() {
	schemaCmd.AddCommand(schema.LoadCmd)
	schemaCmd.AddCommand(schema.DumpCmd)
	RootCmd.AddCommand(schemaCmd)
}
