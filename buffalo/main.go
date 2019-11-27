package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/gobuffalo/buffalo-cli/cli"
)

func main() {
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	b, err := cli.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := b.Main(ctx, os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
