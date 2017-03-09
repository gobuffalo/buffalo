package grifts

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/markbates/grift/grift"
)

var _ = grift.Add("deplist", func(c *grift.Context) error {
	//deplist | grep -v "buffalo" | wc -l
	cmd := exec.Command("deplist", "|", "grep", "-v", "gobuffalo/buffalo")
	out, err := cmd.Output()
	if err != nil {
		return err
	}
	w, err := os.Create("deplist")
	if err != nil {
		return err
	}
	defer w.Close()
	w.Write(bytes.TrimSpace(out))
	return nil
})
