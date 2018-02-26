package updater

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type lockToml struct {
	Name     string   `toml:"name"`
	Branch   string   `toml:"branch"`
	Packages []string `toml:"packages"`
	Revision string   `toml:"revision"`
	Version  string   `toml:"version"`
}

// DepUpdate will attempt to update dependencies belonging to the Buffalo project
func DepUpdate(r *Runner) error {
	if !r.App.WithDep {
		return nil
	}

	if !ask("Would you like to update dependencies related to the Buffalo project?\n(github.com/gobuffalo/* or github.com/markbates/*)") {
		fmt.Println("\tskipping updating dependencies")
		return nil
	}

	p := struct {
		Projects []lockToml `toml:"projects"`
	}{}

	_, err := toml.DecodeFile("./Gopkg.lock", &p)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, l := range p.Projects {
		if strings.Contains(l.Name, "github.com/gobuffalo") || strings.Contains(l.Name, "github.com/markbates") {
			fmt.Printf("~~~ updating Buffalo dependency %s ~~~\n", l.Name)
			cc := exec.Command("dep", "ensure", "-v", "-update", l.Name)
			cc.Stdin = os.Stdin
			cc.Stderr = os.Stderr
			cc.Stdout = os.Stdout
			if err := cc.Run(); err != nil {
				errors.WithStack(err)
			}
		}
	}

	return nil
}

// DepEnsure runs `dep ensure -v` to make sure that any newly changed
// imports are added to dep.
func DepEnsure(r *Runner) error {
	if !r.App.WithDep {
		return nil
	}
	fmt.Println("~~~ Running dep ensure ~~~")
	cc := exec.Command("dep", "ensure", "-v")
	cc.Stdin = os.Stdin
	cc.Stderr = os.Stderr
	cc.Stdout = os.Stdout
	return cc.Run()
}
