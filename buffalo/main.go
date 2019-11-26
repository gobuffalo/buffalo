package main

import (
	"context"
	"log"
	"os"

	"github.com/gobuffalo/buffalo-cli/cli"
)

func main() {
	ctx := context.Background()
	err := cli.Main(ctx, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
