package generate

import (
	"os/exec"

	"github.com/markbates/gentronics"
)

// Fmt is command that will use `goimports` if available,
// or fail back to `gofmt` otherwise.
var Fmt *gentronics.Command

func init() {
	c := "gofmt"
	_, err := exec.LookPath("goimports")
	if err == nil {
		c = "goimports"
	}
	Fmt = gentronics.NewCommand(exec.Command(c, "-w", "."))
}
