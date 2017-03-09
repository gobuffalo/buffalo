package grifts

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/markbates/grift/grift"
)

var depListCmd = exec.Command("deplist", "|", "grep", "-v", "gobuffalo/buffalo")
var _ = grift.Add("deplist", func(c *grift.Context) error {
	out, err := depListCmd.Output()
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

var _ = grift.Add("deplist:count", func(c *grift.Context) error {
	out, err := depListCmd.Output()
	if err != nil {
		return err
	}
	out = bytes.TrimSpace(out)
	l := len(bytes.Split(out, []byte("\n")))
	fmt.Printf("%d Dependencies\n", l)
	return nil
})
