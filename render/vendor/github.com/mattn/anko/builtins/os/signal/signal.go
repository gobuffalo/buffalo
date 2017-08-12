// Package signal implements signal interface for anko script.
package signal

import (
	pkg "os/signal"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("os/signal")

	//m.Define("Ignore", pkg.Ignore)
	m.Define("Notify", pkg.Notify)
	//m.Define("Reset", pkg.Reset)
	m.Define("Stop", pkg.Stop)
	return m
}
