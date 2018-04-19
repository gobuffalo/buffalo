package grifts

import (
	"os"
	"os/exec"

	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
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
		if _, err := exec.LookPath("gometalinter"); err != nil {
			if err := run("go", "get", "-u", "-v", "github.com/alecthomas/gometalinter"); err != nil {
				return errors.WithStack(err)
			}
			if err := run("gometalinter", "--install"); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	var _ = grift.Add("lint", func(c *grift.Context) error {
		if err := grift.Run("tools:install", c); err != nil {
			return err
		}
		return run("gometalinter", "--vendor", "--deadline=3m", "./...")
	})

})
