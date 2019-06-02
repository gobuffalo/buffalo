package grifts

import (
	"os"
	"os/exec"

	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("tools", func() {

	var run = func(args ...string) error {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}

	var _ = grift.Add("install", func(c *grift.Context) error {
		if _, err := exec.LookPath("golangci-lint"); err != nil {
			if err := run("go", "get", "-v", "github.com/golangci/golangci-lint/cmd/golangci-lint"); err != nil {
				return err
			}
		}
		return nil
	})

	var _ = grift.Add("lint", func(c *grift.Context) error {
		if err := grift.Run("tools:install", c); err != nil {
			return err
		}

		return run("golangci-lint", "run", "--fast", "--deadline=3m")
	})

})
