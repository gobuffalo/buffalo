// Package rand implements math/rand interface for anko script.
package rand

import (
	t "math/rand"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("rand")
	m.Define("ExpFloat64", t.ExpFloat64)
	m.Define("Float32", t.Float32)
	m.Define("Float64", t.Float64)
	m.Define("Int", t.Int)
	m.Define("Int31", t.Int31)
	m.Define("Int31n", t.Int31n)
	m.Define("Int63", t.Int63)
	m.Define("Int63n", t.Int63n)
	m.Define("Intn", t.Intn)
	m.Define("NormFloat64", t.NormFloat64)
	m.Define("Perm", t.Perm)
	m.Define("Seed", t.Seed)
	m.Define("Uint32", t.Uint32)
	return m
}
