package cmd

import (
	"os"

	"github.com/markbates/pop"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var migrationPath string

var migrateCmd = &cobra.Command{
	Use:     "migrate",
	Aliases: []string{"m"},
	Short:   "Runs migrations against your database.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		RootCmd.PersistentPreRun(cmd, args)
		return os.MkdirAll(migrationPath, 0766)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mig, err := pop.NewFileMigrator(migrationPath, getConn())
		if err != nil {
			return errors.WithStack(err)
		}
		return mig.Up()
	},
}

func init() {
	RootCmd.AddCommand(migrateCmd)
	RootCmd.PersistentFlags().StringVarP(&migrationPath, "path", "p", "./migrations", "Path to the migrations folder")
}
