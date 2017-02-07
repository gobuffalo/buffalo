package cmd

import "github.com/markbates/pop/soda/cmd"

func init() {
	c := cmd.RootCmd
	c.Use = "db"
	RootCmd.AddCommand(c)
}
