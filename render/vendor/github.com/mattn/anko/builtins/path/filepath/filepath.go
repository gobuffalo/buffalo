// Package path implements path manipulation interface for anko script.
package filepath

import (
	f "path/filepath"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("filepath")
	m.Define("Join", f.Join)
	m.Define("Clean", f.Join)
	m.Define("Abs", f.Abs)
	m.Define("Base", f.Base)
	m.Define("Clean", f.Clean)
	m.Define("Dir", f.Dir)
	m.Define("EvalSymlinks", f.EvalSymlinks)
	m.Define("Ext", f.Ext)
	m.Define("FromSlash", f.FromSlash)
	m.Define("Glob", f.Glob)
	m.Define("HasPrefix", f.HasPrefix)
	m.Define("IsAbs", f.IsAbs)
	m.Define("Join", f.Join)
	m.Define("Match", f.Match)
	m.Define("Rel", f.Rel)
	m.Define("Split", f.Split)
	m.Define("SplitList", f.SplitList)
	m.Define("ToSlash", f.ToSlash)
	m.Define("VolumeName", f.VolumeName)
	return m
}
