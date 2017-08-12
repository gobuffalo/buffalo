// Package json implements json interface for anko script.
package json

import (
	"encoding/json"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("json")
	m.Define("Marshal", json.Marshal)
	m.Define("Unmarshal", json.Unmarshal)
	return m
}
