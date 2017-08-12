// Package path implements path interface for anko script.
package path

import (
	pkg "path"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("path")
	m.Define("Base", pkg.Base)
	m.Define("Clean", pkg.Clean)
	m.Define("Dir", pkg.Dir)
	m.Define("ErrBadPattern", pkg.ErrBadPattern)
	m.Define("Ext", pkg.Ext)
	m.Define("IsAbs", pkg.IsAbs)
	m.Define("Join", pkg.Join)
	m.Define("Match", pkg.Match)
	m.Define("Split", pkg.Split)
	return m
}
