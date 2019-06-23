package fix

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gogen"
)

var rTools = []string{}

func installTools(r *Runner) error {
	fmt.Println("~~~ Installing required tools ~~~")
	run := genny.WetRunner(context.Background())
	g := genny.New()
	app := r.App
	if app.WithDep {
		if _, err := exec.LookPath("dep"); err != nil {
			g.RunFn(gogen.Install("github.com/golang/dep/cmd/dep"))
		}
	}
	if app.WithPop {
		rTools = append(rTools, "github.com/gobuffalo/buffalo-pop")
	}
	for _, t := range rTools {
		g.Command(gogen.Get(t))
	}
	run.With(g)
	return run.Run()
}
