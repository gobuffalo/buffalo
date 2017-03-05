package generators

import (
	"os/exec"
)

// GoInstall compiles and installs packages and dependencies
func GoInstall(pkg string) *exec.Cmd {
	args := []string{"install"}
	args = append(args, pkg)
	return exec.Command("go", args...)
}

// GoGet downloads and installs packages and dependencies
func GoGet(pkg string) *exec.Cmd {
	args := []string{"get"}
	args = append(args, pkg)
	return exec.Command("go", args...)
}

// GoFmt is command that will use `goimports` if available,
// or fail back to `gofmt` otherwise.
func GoFmt() *exec.Cmd {
	c := "gofmt"
	_, err := exec.LookPath("goimports")
	if err == nil {
		c = "goimports"
	}
	_, err = exec.LookPath("gofmt")
	if err != nil {
		return exec.Command("echo", "could not find gofmt or goimports")
	}
	return exec.Command(c, "-w", ".")
}
