// Package exec implements os/exec interface for anko script.
package exec

import (
	e "os/exec"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("exec")
	m.Define("ErrNotFound", e.ErrNotFound)
	m.Define("LookPath", e.LookPath)
	m.Define("Command", e.Command)
	return m
}
