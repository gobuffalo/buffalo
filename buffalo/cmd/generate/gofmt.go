package generate

import (
	"os/exec"

	"github.com/markbates/gentronics"
)

var Fmt *gentronics.Command

func init() {
	c := "gofmt"
	_, err := exec.LookPath("goimports")
	if err == nil {
		c = "goimports"
	}
	Fmt = gentronics.NewCommand(exec.Command(c, "-w", "."))
}
