package cmd

import (
	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Apply all of the 'up' migrations.",
	RunE: func(cmd *cobra.Command, args []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return errors.WithStack(err)
		}
		return mig.Up()
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
}
