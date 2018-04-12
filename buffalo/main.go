package main

import (
	"log"

	"github.com/gobuffalo/buffalo/buffalo/cmd"
	"github.com/gobuffalo/buffalo/internal/vbuffalo"
)

func main() {
	err := vbuffalo.Execute(func() error {
		cmd.Execute()
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
