package grifts

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/markbates/grift/grift"
)

var _ = grift.Add("completions:fish", func(c *grift.Context) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	copyCmd := "cp " + filepath.Join(pwd, "completions/fish/buffalo.fish") + " $fish_complete_path[1]"
	cmd := exec.Command("fish", "-c", copyCmd)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
})
