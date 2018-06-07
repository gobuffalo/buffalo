package vbuffalo

import (
	"os"

	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func writeMain() error {
	f, err := os.Create(mainPath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	s, err := plush.Render(mainTemplate, plush.NewContextWith(map[string]interface{}{
		"app":     app,
		"version": runtime.Version,
	}))
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = f.WriteString(s)
	return err
}

const mainTemplate = `package main

import (
	"fmt"
	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/buffalo/buffalo/cmd"
)

func main() {
	fmt.Printf("%s [<%= version %>]\n\n", runtime.Version)
	cmd.Execute()
}
`
