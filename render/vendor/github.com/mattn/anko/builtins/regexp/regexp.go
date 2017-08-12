// Package regexp implements regexp interface for anko script.
package sort

import (
	r "regexp"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("sort")
	m.Define("Match", r.Match)
	m.Define("MatchReader", r.MatchReader)
	m.Define("MatchString", r.MatchString)
	m.Define("QuoteMeta", r.QuoteMeta)
	m.Define("Compile", r.Compile)
	m.Define("CompilePOSIX", r.CompilePOSIX)
	m.Define("MustCompile", r.MustCompile)
	m.Define("MustCompilePOSIX", r.MustCompilePOSIX)
	return m
}
