// Package big implements math/big interface for anko script.
package big

import (
	t "math/big"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewModule("big")
	m.Define("Above", t.Above)
	m.Define("AwayFromZero", t.AwayFromZero)
	m.Define("Below", t.Below)
	m.Define("Exact", t.Exact)
	m.Define("Jacobi", t.Jacobi)
	m.Define("MaxBase", t.MaxBase)
	m.Define("MaxExp", t.MaxExp)
	// TODO https://github.com/mattn/anko/issues/49
	//m.Define("MaxPrec", t.MaxPrec)
	m.Define("MinExp", t.MinExp)
	m.Define("NewFloat", t.NewFloat)
	m.Define("NewInt", t.NewInt)
	m.Define("NewRat", t.NewRat)
	m.Define("ParseFloat", t.ParseFloat)
	m.Define("ToNearestAway", t.ToNearestAway)
	m.Define("ToNearestEven", t.ToNearestEven)
	m.Define("ToNegativeInf", t.ToNegativeInf)
	m.Define("ToPositiveInf", t.ToPositiveInf)
	m.Define("ToZero", t.ToZero)
	return m
}
