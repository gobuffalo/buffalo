package fix

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
)

var rTools = []string{}

func installTools(r *Runner) error {
	fmt.Println("~~~ Installing required tools ~~~")
	run := genny.WetRunner(context.Background())
	app := r.App
	if app.WithDep {
		if _, err := exec.LookPath("dep"); err != nil {
			run.WithRun(gotools.Install("github.com/golang/dep/cmd/dep"))
		}
	}
	if app.WithPop {
		rTools = append(rTools, "github.com/gobuffalo/buffalo-pop")
	}
	for _, t := range rTools {
		run.WithRun(gotools.Get(t))
	}
	return run.Run()
}
