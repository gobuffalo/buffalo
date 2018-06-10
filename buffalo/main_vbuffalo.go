// +build vbuffalo

package main

import (
	"fmt"
	"log"

	"github.com/gobuffalo/buffalo/buffalo/cmd"
	"github.com/gobuffalo/buffalo/internal/vbuffalo"
)

func main() {
	fmt.Println("vbuffalo")
	err := vbuffalo.Execute(func() error {
		cmd.Execute()
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
