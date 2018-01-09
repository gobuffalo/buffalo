package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
	cmd := newTestCmd(args)
	var runFlag bool
	var mFlag bool
	for i, a := range args {
		if a == "-run" {
			runFlag = true
		}
		if a == "-m" {
			mFlag = true
			args[i] = "-testify.m"
		}
	}

	if mFlag {
		return mFlagRunner(args)
	}

	if !runFlag {
		pkgs, err := testPackages()
		if err != nil {
			return errors.WithStack(err)
		}
		cmd.Args = append(cmd.Args, pkgs...)
	}
	logrus.Info(strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func mFlagRunner(args []string) error {
	app := meta.New(".")
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)

	pkgs, err := testPackages()
	if err != nil {
		return errors.WithStack(err)
	}
	var errs bool
	for _, p := range pkgs {
		os.Chdir(pwd)
		if p == app.PackagePkg {
			continue
		}
		cmd := newTestCmd(args)
		p = strings.TrimPrefix(p, app.PackagePkg+string(filepath.Separator))
		logrus.Info(strings.Join(cmd.Args, " "))
		os.Chdir(p)
		if err := cmd.Run(); err != nil {
			errs = true
		}
	}
	if errs {
		return errors.New("errors running tests")
	}
	return nil
}

func testPackages() ([]string, error) {
	args := []string{}
	out, err := exec.Command(envy.Get("GO_BIN", "go"), "list", "./...").Output()
	if err != nil {
		return args, err
	}
	pkgs := bytes.Split(bytes.TrimSpace(out), []byte("\n"))
	for _, p := range pkgs {
		if !vendorRegex.Match(p) {
			args = append(args, string(p))
		}
	}
	return args, nil
}

func newTestCmd(args []string) *exec.Cmd {
	cmd := exec.Command(envy.Get("GO_BIN", "go"), "test", "-p", "1")
	if _, err := exec.LookPath("gotest"); err == nil {
		cmd = exec.Command("gotest", "-p", "1")
	}
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
