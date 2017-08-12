// Package fmt implements json interface for anko script.
package fmt

import (
	pkg "fmt"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("fmt")
	m.Define("Errorf", pkg.Errorf)
	m.Define("Fprint", pkg.Fprint)
	m.Define("Fprintf", pkg.Fprintf)
	m.Define("Fprintln", pkg.Fprintln)
	m.Define("Fscan", pkg.Fscan)
	m.Define("Fscanf", pkg.Fscanf)
	m.Define("Fscanln", pkg.Fscanln)
	m.Define("Print", pkg.Print)
	m.Define("Printf", pkg.Printf)
	m.Define("Println", pkg.Println)
	m.Define("Scan", pkg.Scan)
	m.Define("Scanf", pkg.Scanf)
	m.Define("Scanln", pkg.Scanln)
	m.Define("Sprint", pkg.Sprint)
	m.Define("Sprintf", pkg.Sprintf)
	m.Define("Sprintln", pkg.Sprintln)
	m.Define("Sscan", pkg.Sscan)
	m.Define("Sscanf", pkg.Sscanf)
	m.Define("Sscanln", pkg.Sscanln)
	return m
}
