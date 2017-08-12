package generate

import (
	"github.com/pkg/errors"

	"github.com/markbates/going/defaults"
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var FizzCmd = &cobra.Command{
	Use:     "fizz [name]",
	Aliases: []string{"migration"},
	Short:   "Generates Up/Down migrations for your database using fizz.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("You must supply a name for your migration!")
		}
		cflag := cmd.Flag("path")
		migrationPath := defaults.String(cflag.Value.String(), "./migrations")
		return pop.MigrationCreate(migrationPath, args[0], "fizz", nil, nil)
	},
}
