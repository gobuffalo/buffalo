package updater

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/gobuffalo/envy"
	"github.com/markbates/deplist"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type lockToml struct {
	Name     string   `toml:"name"`
	Branch   string   `toml:"branch"`
	Packages []string `toml:"packages"`
	Revision string   `toml:"revision"`
	Version  string   `toml:"version"`
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

// DepEnsure runs `dep ensure -v` to make sure that any newly changed
// imports are added to dep.
func DepEnsure(r *Runner) error {
	if !r.App.WithDep {
		return goGetUpdate(r)
	}
	fmt.Println("~~~ Running dep ensure ~~~")
	cc := exec.Command("dep", "ensure", "-v")
	cc.Stdin = os.Stdin
	cc.Stderr = os.Stderr
	cc.Stdout = os.Stdout
	if err := cc.Run(); err != nil {
		return errors.WithStack(err)
	}

	apkg := []string{
		"github.com/gobuffalo/tags@v2.0.0",
		"github.com/gobuffalo/suite@v2.0.0",
	}
	args := []string{"ensure", "-v", "-add"}

	for _, p := range apkg {
		args = append(args, p)
	}
	cc = exec.Command("dep", args...)
	cc.Stdin = os.Stdin
	cc.Stderr = os.Stderr
	cc.Stdout = os.Stdout
	if err := cc.Run(); err != nil {
		return errors.WithStack(err)
	}

	upkg := []string{
		"github.com/gobuffalo/buffalo",
		"github.com/gobuffalo/plush",
		"github.com/markbates/inflect",
	}

	args = []string{"ensure", "-v", "-update"}
	for _, p := range upkg {
		args = append(args, p)
	}
	cc = exec.Command("dep", args...)
	cc.Stdin = os.Stdin
	cc.Stderr = os.Stderr
	cc.Stdout = os.Stdout
	return cc.Run()
}
