// +build !appengine

// Package url implements url interface for anko script.
package url

import (
	u "net/url"

	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	m := env.NewPackage("url")
	m.DefineType("Values", make(u.Values))
	m.Define("Parse", u.Parse)
	return m
}
