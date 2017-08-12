package main

import (
	"fmt"
	"log"

	"github.com/mattn/anko/vm"
)

func main() {
	env := vm.NewEnv()

	env.Define("foo", 1)
	env.Define("bar", func() int {
		return 2
	})

	v, err := env.Execute(`foo + bar()`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v)
}
