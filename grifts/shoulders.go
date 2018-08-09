package grifts

import (
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/shoulders/shoulders"
	"github.com/markbates/grift/grift"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var _ = grift.Desc("shoulders", "Prints a listing all of the 3rd party packages used by buffalo.")
var _ = grift.Add("shoulders:list", func(c *grift.Context) error {
	view, err := shoulders.New()
	if err != nil {
		return errors.WithStack(err)
	}
	logrus.Infof(strings.Join(view.Deps, "\n"))
	return nil
})

var _ = grift.Desc("shoulders", "Generates a file listing all of the 3rd party packages used by buffalo.")
var _ = grift.Add("shoulders", func(c *grift.Context) error {
	view, err := shoulders.New()
	view.Name = "Buffalo"
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.Create(path.Join(envy.GoPath(), "src", "github.com", "gobuffalo", "buffalo", "SHOULDERS.md"))
	if err != nil {
		return err
	}

	if err := view.Write(f); err != nil {
		return err
	}

	return commitAndPushShoulders()
})

func commitAndPushShoulders() error {
	cmd := exec.Command("git", "commit", "SHOULDERS.md", "-m", "Updated SHOULDERS.md")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "push", "origin")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
