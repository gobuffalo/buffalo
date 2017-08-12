package cmd

import (
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var migrateResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "The equivalent of running `migrate down` and then `migrate up`",
	RunE: func(cmd *cobra.Command, args []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return errors.WithStack(err)
		}
		return mig.Reset()
	},
}

func init() {
	migrateCmd.AddCommand(migrateResetCmd)
}
