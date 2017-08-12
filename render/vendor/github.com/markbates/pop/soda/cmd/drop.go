package cmd

import (
	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

var all bool

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drops databases for you",
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			for _, conn := range pop.Connections {
				pop.DropDB(conn)
			}
		} else {
			pop.DropDB(getConn())
		}
	},
}

func init() {
	dropCmd.Flags().BoolVarP(&all, "all", "a", false, "Drops all of the databases in the database.yml")
	RootCmd.AddCommand(dropCmd)
}
