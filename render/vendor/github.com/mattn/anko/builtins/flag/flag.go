// Package flag implements flag interface for anko script.
package flag

import (
	pkg "flag"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("flag")
	m.Define("Arg", pkg.Arg)
	m.Define("Args", pkg.Args)
	m.Define("Bool", pkg.Bool)
	m.Define("BoolVar", pkg.BoolVar)
	m.Define("CommandLine", pkg.CommandLine)
	m.Define("ContinueOnError", pkg.ContinueOnError)
	m.Define("Duration", pkg.Duration)
	m.Define("DurationVar", pkg.DurationVar)
	m.Define("ErrHelp", pkg.ErrHelp)
	m.Define("ExitOnError", pkg.ExitOnError)
	m.Define("Float64", pkg.Float64)
	m.Define("Float64Var", pkg.Float64Var)
	m.Define("Int", pkg.Int)
	m.Define("Int64", pkg.Int64)
	m.Define("Int64Var", pkg.Int64Var)
	m.Define("IntVar", pkg.IntVar)
	m.Define("Lookup", pkg.Lookup)
	m.Define("NArg", pkg.NArg)
	m.Define("NFlag", pkg.NFlag)
	m.Define("NewFlagSet", pkg.NewFlagSet)
	m.Define("PanicOnError", pkg.PanicOnError)
	m.Define("Parse", pkg.Parse)
	m.Define("Parsed", pkg.Parsed)
	m.Define("PrintDefaults", pkg.PrintDefaults)
	m.Define("Set", pkg.Set)
	m.Define("String", pkg.String)
	m.Define("StringVar", pkg.StringVar)
	m.Define("Uint", pkg.Uint)
	m.Define("Uint64", pkg.Uint64)
	m.Define("Uint64Var", pkg.Uint64Var)
	m.Define("UintVar", pkg.UintVar)
	m.Define("Usage", pkg.Usage)
	m.Define("Var", pkg.Var)
	m.Define("Visit", pkg.Visit)
	m.Define("VisitAll", pkg.VisitAll)
	return m
}
