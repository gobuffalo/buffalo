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

	"github.com/gobuffalo/pop"
	"github.com/spf13/cobra"
)

const vendorPattern = "/vendor/"

var vendorRegex = regexp.MustCompile(vendorPattern)
var forceMigrations = false

func init() {
	decorate("test", testCmd)
	testCmd.Flags().BoolVarP(&forceMigrations, "force-migrations", "m", false, "skips loading the schema and instead runs migrations before running tests")

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

			if forceMigrations {
				fm, err := pop.NewFileMigrator("./migrations", test)

				if err != nil {
					return err
				}

				if err := fm.Up(); err != nil {
					return err
				}

				return testRunner(args)
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
		fm, err := pop.NewFileMigrator("./migrations", test)
		if err != nil {
			return nil
		}

		if err := fm.Up(); err == nil {
			if f, err := os.Open(filepath.Join("migrations", "schema.sql")); err == nil {
				return f
			}
		}
	}
	return nil
}

func testRunner(args []string) error {
	var mFlag bool
	var query string
	cargs := []string{}
	pargs := []string{}
	var larg string
	for i, a := range args {
		switch a {
		case "-run", "-m":
			query = args[i+1]
			mFlag = true
		case "-v":
			cargs = append(cargs, "-v")
		default:
			if larg != "-run" && larg != "-m" {
				pargs = append(pargs, a)
			}
		}
		larg = a
	}

	cmd := newTestCmd(cargs)
	if mFlag {
		return mFlagRunner{
			query: query,
			args:  cargs,
			pargs: pargs,
		}.Run()
	}

	pkgs, err := testPackages(pargs)
	if err != nil {
		return errors.WithStack(err)
	}
	cmd.Args = append(cmd.Args, pkgs...)
	logrus.Info(strings.Join(cmd.Args, " "))
	return cmd.Run()
}

type mFlagRunner struct {
	query string
	args  []string
	pargs []string
}

func (m mFlagRunner) Run() error {
	app := meta.New(".")
	pwd, _ := os.Getwd()
	defer os.Chdir(pwd)

	pkgs, err := testPackages(m.pargs)
	if err != nil {
		return errors.WithStack(err)
	}
	var errs bool
	for _, p := range pkgs {
		os.Chdir(pwd)
		if p == app.PackagePkg {
			continue
		}
		p = strings.TrimPrefix(p, app.PackagePkg+string(filepath.Separator))
		os.Chdir(p)

		cmd := newTestCmd(m.args)
		if hasTestify(p) {
			cmd.Args = append(cmd.Args, "-testify.m", m.query)
		} else {
			cmd.Args = append(cmd.Args, "-run", m.query)
		}
		logrus.Info(strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			errs = true
		}
	}
	if errs {
		return errors.New("errors running tests")
	}
	return nil
}

func hasTestify(p string) bool {
	cmd := exec.Command("go", "test", "-thisflagdoesntexist")
	b, _ := cmd.Output()
	return bytes.Contains(b, []byte("-testify.m"))
}

func testPackages(givenArgs []string) ([]string, error) {
	// If there are args, then assume these are the packages to test.
	//
	// Instead of always returning all packages from 'go list ./...', just
	// return the given packages in this case
	if len(givenArgs) > 0 {
		return givenArgs, nil
	}
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
	cargs := []string{"test", "-p", "1"}
	app := meta.New(".")
	tf := `-tags "` + app.BuildTags("development").String()
	cargs = append(cargs, tf)
	cargs = append(cargs, args...)
	cmd := exec.Command(envy.Get("GO_BIN", "go"), cargs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}
