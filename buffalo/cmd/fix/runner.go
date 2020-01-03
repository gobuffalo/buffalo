package fix

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo/runtime"
	"github.com/gobuffalo/meta"
)

// Check interface for runnable checker functions
type Check func(*Runner) error

// Runner will run all compatible checks
type Runner struct {
	App      meta.App
	Warnings []string
}

// Run all compatible checks
func Run() error {
	fmt.Printf("! This updater will attempt to update your application to Buffalo version: %s\n", runtime.Version)
	if !ask("Do you wish to continue?") {
		fmt.Println("~~~ cancelling update ~~~")
		return nil
	}

	r := &Runner{
		App:      meta.New("."),
		Warnings: []string{},
	}

	defer func() {
		if len(r.Warnings) == 0 {
			return
		}

		fmt.Println("\n\n----------------------------")
		fmt.Printf("!!! (%d) Warnings Were Found !!!\n\n", len(r.Warnings))
		for _, w := range r.Warnings {
			fmt.Printf("[WARNING]: %s\n", w)
		}
	}()

	for _, c := range checks {

		cmd := exec.Command("go", "mod", "tidy")
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}

		if err := c(r); err != nil {
			return err
		}
	}
	return nil
}
