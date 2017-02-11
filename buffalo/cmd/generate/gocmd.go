package generate

import "os/exec"

// GoInstall compiles and installs packages and dependencies
func GoInstall(pkg string) *exec.Cmd {
	args := []string{"install"}
	args = append(args, pkg)
	return exec.Command("go", args...)
}

// GoGet downloads and installs packages and dependencies
func GoGet(pkg string, buildFlags ...string) *exec.Cmd {
	args := []string{"get"}
	args = append(args, pkg)
	args = append(args, buildFlags...)
	return exec.Command("go", args...)
}
