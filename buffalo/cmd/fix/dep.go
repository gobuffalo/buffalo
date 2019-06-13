package fix

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
)

// packages to add to Gopkg.toml
var apkg = []string{}

// packages ensure get updated
var upkg = []string{
	"github.com/gobuffalo/buffalo",
	"github.com/gobuffalo/plush",
	"github.com/gobuffalo/events",
	"github.com/gobuffalo/suite",
	"github.com/gobuffalo/flect",
}

// DepEnsure runs `dep ensure -v` or `go get -u` depending on app tooling
// to make sure that any newly changed imports are added to dep or installed.
func DepEnsure(r *Runner) error {
	if r.App.WithPop {
		upkg = append(upkg, "github.com/gobuffalo/fizz", "github.com/gobuffalo/pop")
	}
	if !r.App.WithDep {
		fmt.Println("~~~ Running go get ~~~")
		return modGetUpdate(r)
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
				return err
			}
		}
	}

	if len(upkg) > 0 {
		args := []string{"ensure", "-v", "-update"}
		if err := depRunner(args); err != nil {
			return err
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
	run := genny.WetRunner(context.Background())
	g := genny.New()
	for _, x := range upkg {
		if x == "github.com/gobuffalo/buffalo" {
			continue
		}
		g.Command(gogen.Get(x))
	}

	for _, x := range []string{"beta", "rc"} {
		if !strings.Contains(runtime.Version, x) {
			continue
		}
		g.Command(gogen.Get("github.com/gobuffalo/buffalo@"+runtime.Version, "-u"))
	}
	run.With(g)
	return run.Run()
}
