package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"

	"github.com/markbates/pop"
	"github.com/spf13/cobra"
)

const vendorPattern = "/vendor/"

var vendorRegex = regexp.MustCompile(vendorPattern)

func init() {
	decorate("test", testCmd)
	RootCmd.AddCommand(testCmd)
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
				return errors.WithStack(err)
			}

			// drop the test db:
			test.Dialect.DropDB()

			// create the test db:
			err = test.Dialect.CreateDB()
			if err != nil {
				return errors.WithStack(err)
			}

			if schema := findSchema(); schema != nil {
				err = test.Dialect.LoadSchema(schema)
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
		return testRunner(args)
	},
}

func findSchema() io.Reader {
	if f, err := os.Open(filepath.Join("migrations", "schema.sql")); err == nil {
		return f
	}
	if dev, err := pop.Connect("development"); err == nil {
		schema := &bytes.Buffer{}
		if err = dev.Dialect.DumpSchema(schema); err == nil {
			return schema
		}
	}

	if test, err := pop.Connect("test"); err == nil {
		if err := test.MigrateUp("./migrations"); err == nil {
			if f, err := os.Open(filepath.Join("migrations", "schema.sql")); err == nil {
				return f
			}
		}
	}
	return nil
}

func testRunner(args []string) error {
	cmd := exec.Command(envy.Get("GO_BIN", "go"), "test", "-p", "1")
	if _, err := exec.LookPath("gotest"); err == nil {
		cmd = exec.Command("gotest", "-p", "1")
	}
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
		out, err := exec.Command(envy.Get("GO_BIN", "go"), "list", "./...").Output()
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
