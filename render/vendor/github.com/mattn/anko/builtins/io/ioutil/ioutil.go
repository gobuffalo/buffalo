// Package ioutil implements I/O interface for anko script.
package ioutil

import (
	u "io/ioutil"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("iotuil")
	m.Define("ReadAll", u.ReadAll)
	m.Define("ReadDir", u.ReadDir)
	m.Define("ReadFile", u.ReadFile)
	m.Define("WriteFile", u.WriteFile)
	return m
}
