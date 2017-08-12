package schema

import (
	"io"
	"os"
	"path/filepath"

	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var dumpOptions = struct {
	env    string
	output string
}{}

var DumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dumps out the schema of the selected database",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := pop.Connect(dumpOptions.env)
		if err != nil {
			return err
		}
		var out io.Writer
		if dumpOptions.output == "-" {
			out = os.Stdout
		} else {
			err = os.MkdirAll(filepath.Dir(dumpOptions.output), 0755)
			if err != nil {
				return err
			}
			out, err = os.Create(dumpOptions.output)
			if err != nil {
				return err
			}
		}
		err = c.Dialect.DumpSchema(out)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	DumpCmd.Flags().StringVarP(&dumpOptions.env, "env", "e", "development", "The environment you want to run schema against. Will use $GO_ENV if set.")
	DumpCmd.Flags().StringVarP(&dumpOptions.output, "output", "o", "schema.sql", "The path to dump the schema to.")
}
