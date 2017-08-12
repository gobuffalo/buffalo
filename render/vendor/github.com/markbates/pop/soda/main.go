package main

import "github.com/markbates/pop/soda/cmd"

func main() {
	cmd.RootCmd.Use = "soda"
	cmd.Execute()
}
