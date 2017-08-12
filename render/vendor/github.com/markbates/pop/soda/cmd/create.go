package cmd

import (
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates databases for you",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if all {
			for _, conn := range pop.Connections {
				err = pop.CreateDB(conn)
				if err != nil {
					return err
				}
			}
		} else {
			err = pop.CreateDB(getConn())
		}
		return err
	},
}

func init() {
	createCmd.Flags().BoolVarP(&all, "all", "a", false, "Creates all of the databases in the database.yml")
	RootCmd.AddCommand(createCmd)
}
