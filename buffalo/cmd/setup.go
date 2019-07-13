package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/events"
	"github.com/gobuffalo/meta"
	"github.com/markbates/deplist"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var setupOptions = struct {
	verbose       bool
	updateGoDeps  bool
	dropDatabases bool
}{}

type setupCheck func(meta.App) error

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup a newly created, or recently checked out application.",
	Long: `Setup runs through checklist to make sure dependencies are setup correctly.

Dependencies (if used):
* Runs "dep ensure" to install required Go dependencies.

Asset Pipeline (if used):
* Runs "npm install" or "yarn install" to install asset dependencies.

Database (if used):
* Runs "buffalo db create -a" to create databases.
* Runs "buffalo db migrate" to run database migrations.
* Runs "buffalo task db:seed" to seed the database (if the task exists).

Tests:
* Runs "buffalo test" to confirm the application's tests are running properly.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		app := meta.New(".")
		payload := events.Payload{
			"app": app,
		}
		events.EmitPayload(EvtSetupStarted, payload)
		for _, check := range []setupCheck{assetCheck, updateGoDepsCheck, databaseCheck, testCheck} {
			err := check(app)
			if err != nil {
				events.EmitError(EvtSetupErr, err, payload)
				return err
			}
		}
		events.EmitPayload(EvtSetupFinished, payload)
		return nil
	},
}

func updateGoDepsCheck(app meta.App) error {
	if app.WithModules {
		c := exec.Command(envy.Get("GO_BIN", "go"), "get")
		return run(c)
	}
	if app.WithDep {
		if _, err := exec.LookPath("dep"); err != nil {
			if err := run(exec.Command(envy.Get("GO_BIN", "go"), "get", "github.com/golang/dep/cmd/dep")); err != nil {
				return err
			}
		}
		args := []string{"ensure"}
		if setupOptions.verbose {
			args = append(args, "-v")
		}
		if setupOptions.updateGoDeps {
			args = append(args, "--update")
		}
		err := run(exec.Command("dep", args...))
		if err != nil {
			return err
		}
		return nil
	}

	// go old school with the installation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg, _ := errgroup.WithContext(ctx)
	deps, err := deplist.List()
	if err != nil {
		return err
	}

	deps["github.com/gobuffalo/suite"] = "github.com/gobuffalo/suite"

	for dep := range deps {
		args := []string{"get"}
		if setupOptions.verbose {
			args = append(args, "-v")
		}
		if setupOptions.updateGoDeps {
			args = append(args, "-u")
		}
		args = append(args, dep)
		c := exec.Command(envy.Get("GO_BIN", "go"), args...)
		f := func() error {
			return run(c)
		}
		wg.Go(f)
	}
	err = wg.Wait()
	if err != nil {
		return fmt.Errorf("We encountered the following error trying to install and update the dependencies for this application:\n%s", err)
	}
	return nil
}

func testCheck(meta.App) error {
	err := run(exec.Command("buffalo", "test"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to run your applications tests:\n%s", err)
	}
	return nil
}

func databaseCheck(app meta.App) error {
	if !app.WithPop {
		return nil
	}
	for _, check := range []setupCheck{dbCreateCheck, dbMigrateCheck, dbSeedCheck} {
		err := check(app)
		if err != nil {
			return err
		}
	}
	return nil
}

func dbCreateCheck(meta.App) error {
	if setupOptions.dropDatabases {
		err := run(exec.Command("buffalo", "pop", "drop", "-a"))
		if err != nil {
			return fmt.Errorf("We encountered an error when trying to drop your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
		}
	}
	err := run(exec.Command("buffalo", "pop", "create", "-a"))
	if err != nil {
		return fmt.Errorf("We encountered an error when trying to create your application's databases. Please check to make sure that your database server is running and that the username and passwords found in the database.yml are properly configured and set up on your database server.\n %s", err)
	}
	return nil
}

func dbMigrateCheck(meta.App) error {
	err := run(exec.Command("buffalo", "pop", "migrate"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to migrate your database:\n%s", err)
	}
	return nil
}

func dbSeedCheck(meta.App) error {
	cmd := exec.Command("buffalo", "t", "list")
	out, err := cmd.Output()
	if err != nil {
		// no tasks configured, so return
		return nil
	}
	if bytes.Contains(out, []byte("db:seed")) {
		err := run(exec.Command("buffalo", "task", "db:seed"))
		if err != nil {
			return fmt.Errorf("We encountered the following error when trying to seed your database:\n%s", err)
		}
	}
	return nil
}

func assetCheck(app meta.App) error {
	if !app.WithWebpack {
		return nil
	}
	if app.WithYarn {
		return yarnCheck(app)
	}
	return npmCheck(app)
}

func npmCheck(app meta.App) error {
	err := nodeCheck(app)
	if err != nil {
		return err
	}
	err = run(exec.Command("npm", "install", "--no-progress"))
	if err != nil {
		return fmt.Errorf("We encountered the following error when trying to install your asset dependencies using npm:\n%s", err)
	}
	return nil
}

func yarnCheck(app meta.App) error {
	if err := nodeCheck(app); err != nil {
		return err
	}
	if _, err := exec.LookPath("yarnpkg"); err != nil {
		err := run(exec.Command("npm", "install", "-g", "yarn"))
		if err != nil {
			return fmt.Errorf("This application require yarn, and we could not find it installed on your system. We tried to install it for you, but ran into the following error:\n%s", err)
		}
	}
	if err := run(exec.Command("yarnpkg", "install", "--no-progress")); err != nil {
		return fmt.Errorf("We encountered the following error when trying to install your asset dependencies using yarn:\n%s", err)
	}
	return nil
}

func nodeCheck(meta.App) error {
	if _, err := exec.LookPath("node"); err != nil {
		return fmt.Errorf("this application requires node, and we could not find it installed on your system please install node and try again")
	}
	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("this application requires npm, and we could not find it installed on your system please install npm and try again")
	}
	return nil
}

func run(cmd *exec.Cmd) error {
	logrus.Infof("--> %s", strings.Join(cmd.Args, " "))
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func init() {
	setupCmd.Flags().BoolVarP(&setupOptions.verbose, "verbose", "v", false, "run with verbose output")
	setupCmd.Flags().BoolVarP(&setupOptions.updateGoDeps, "update", "u", false, "run go get -u against the application's Go dependencies")
	setupCmd.Flags().BoolVarP(&setupOptions.dropDatabases, "drop", "d", false, "drop existing databases")

	decorate("setup", setupCmd)
	RootCmd.AddCommand(setupCmd)
}
