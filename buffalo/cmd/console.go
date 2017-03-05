package cmd

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gobuffalo/makr"
	"github.com/markbates/inflect"
	"github.com/spf13/cobra"
)

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
	Use:     "console",
	Aliases: []string{"c"},
	Short:   "Runs your Buffalo app in a REPL console",
	RunE: func(c *cobra.Command, args []string) error {
		_, err := exec.LookPath("gore")
		if err != nil {
			return errors.New("we could not find \"gore\" in your path. You must first install \"gore\" in order to use the Buffalo console:\n\n$ go get -u github.com/motemen/gore")
		}
		packagePath := packagePath(rootPath)
		packages := []string{}
		for _, p := range []string{"models", "actions"} {
			s, _ := os.Stat(filepath.Join(rootPath, p))
			if s != nil {
				packages = append(packages, filepath.Join(packagePath, p))
			}
		}

		fname := inflect.Parameterize(packagePath) + "_loader.go"
		g := makr.New()
		g.Add(makr.NewFile(fname, cMain))
		err = g.Run(os.TempDir(), makr.Data{
			"packages": packages,
		})
		os.Chdir(rootPath)
		if err != nil {
			return err
		}

		cmd := exec.Command("gore", "-autoimport", "-context", filepath.Join(os.TempDir(), fname))
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()
	},
}

func init() {
	RootCmd.AddCommand(consoleCmd)
}

var cMain = `
package main

{{range .packages}}
import _ "{{.}}"
{{end}}
`
