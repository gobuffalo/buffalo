package fix

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/envy"
	"github.com/markbates/deplist"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// packages to add to Gopkg.toml
var apkg = []string{}

// packages ensure get updated
var upkg = []string{
	"github.com/gobuffalo/buffalo",
	"github.com/gobuffalo/plush",
	"github.com/gobuffalo/pop",
	"github.com/gobuffalo/suite",
	"github.com/markbates/inflect",
}

// DepEnsure runs `dep ensure -v` or `go get -u` depending on app tooling
// to make sure that any newly changed imports are added to dep or installed.
func DepEnsure(r *Runner) error {
	fmt.Println("Running Dependencies")
	if r.App.WithModules {
		fmt.Println("~~~ Updating modules ~~~")
		return modGetUpdate(r)
	}

	if !r.App.WithDep {
		fmt.Println("~~~ Running go get ~~~")
		return goGetUpdate(r)
	}

	fmt.Println("~~~ Running dep ensure ~~~")
	return runDepEnsure(r)
}

func runDepEnsure(r *Runner) error {
	for _, x := range []string{"beta", "rc", "development"} {
		if strings.Contains(runtime.Version, x) {
			r.Warnings = append(r.Warnings, fmt.Sprintf("This is not an official release and you will need to MANUALLY adjust your Gopkg.toml file to use this release."))
			break
		}
	}

	if len(apkg) > 0 {
		args := []string{"ensure", "-v", "-add"}
		args = append(args, apkg...)
		if err := depRunner(args); err != nil {
			// *sigh* - yeah, i know
			if !strings.Contains(err.Error(), "is already in Gopkg.toml") {
				return errors.WithStack(err)
			}
		}
	}

	if len(upkg) > 0 {
		args := []string{"ensure", "-v", "-update"}
		args = append(args, upkg...)
		if err := depRunner(args); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func depRunner(args []string) error {
	cc := exec.Command("dep", args...)
	cc.Stdin = os.Stdin
	cc.Stderr = os.Stderr
	cc.Stdout = os.Stdout
	return cc.Run()
}

func modGetUpdate(r *Runner) error {
	cmd := exec.Command(envy.Get("GO_BIN", "go"), "get", "-u")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return errors.WithStack(err)
	}

	for _, x := range []string{"beta", "rc"} {
		if !strings.Contains(runtime.Version, x) {
			continue
		}
		cmd = exec.Command(envy.Get("GO_BIN", "go"), "get", "-u", "github.com/gobuffalo/buffalo@"+runtime.Version)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
	return nil
}

func goGetUpdate(r *Runner) error {
	fmt.Println("~~~ Running go get ~~~")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg, _ := errgroup.WithContext(ctx)
	deps, err := deplist.List()
	if err != nil {
		return errors.WithStack(err)
	}
	for dep := range deps {
		args := []string{"get", "-u"}
		args = append(args, dep)
		cc := exec.Command(envy.Get("GO_BIN", "go"), args...)
		f := func() error {
			cc.Stdin = os.Stdin
			cc.Stderr = os.Stderr
			cc.Stdout = os.Stdout
			return cc.Run()
		}
		wg.Go(f)
	}
	err = wg.Wait()
	if err != nil {
		return errors.Errorf("We encountered the following error trying to install and update the dependencies for this application:\n%s", err)
	}
	return nil
}
