package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

const vendorPattern = "/vendor/"

var vendorRegex *regexp.Regexp

func init() {
	RootCmd.AddCommand(testCmd)
	vendorRegex = regexp.MustCompile(vendorPattern)
}

var testCmd = &cobra.Command{
	Use:                "test",
	Short:              "Runs the tests for your Buffalo app",
	DisableFlagParsing: true,
	RunE: func(c *cobra.Command, args []string) error {
		os.Setenv("GO_ENV", "test")
		if _, err := os.Stat("database.yml"); err == nil {
			// there's a database
			test, err := pop.Connect("test")
			if err != nil {
				return err
			}

			// drop the test db:
			test.Dialect.DropDB()

			// create the test db:
			err = test.Dialect.CreateDB()
			if err != nil {
				return err
			}

			dev, err := pop.Connect("development")
			if err != nil {
				return err
			}
			schema := &bytes.Buffer{}
			err = dev.Dialect.DumpSchema(schema)
			if err != nil {
				return err
			}

			err = test.Dialect.LoadSchema(schema)
			if err != nil {
				return err
			}
		}
		return testRunner(args)
	},
}

func testRunner(args []string) error {
	cmd := exec.Command("go", "test")
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	runFlag := false
	for _, a := range cmd.Args {
		if a == "-run" {
			runFlag = true
		}
	}
	if !runFlag {
		out, err := exec.Command("go", "list", "./...").Output()
		if err != nil {
			return err
		}
		pkgs := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
		for _, p := range pkgs {
			if !vendorRegex.Match(p) {
				cmd.Args = append(cmd.Args, string(p))
			}
		}
	}
	fmt.Println(strings.Join(cmd.Args, " "))
	return cmd.Run()
}
