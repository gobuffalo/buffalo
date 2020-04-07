package fix

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/gobuffalo/genny/v2"
)

var rTools = []string{}

func installTools(r *Runner) error {
	fmt.Println("~~~ Installing required tools ~~~")
	run := genny.WetRunner(context.Background())
	g := genny.New()
	app := r.App
	if app.WithPop {
		rTools = append(rTools, "github.com/gobuffalo/buffalo-pop/v2")
	}
	for _, t := range rTools {
		g.Command(exec.Command("go", "get", t))
	}
	run.With(g)
	return run.Run()
}
