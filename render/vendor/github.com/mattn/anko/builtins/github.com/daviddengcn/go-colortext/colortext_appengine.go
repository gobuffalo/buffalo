// +build appengine

// Package colortext implements terminal interface for anko script.
package colortext

import (
	"github.com/mattn/anko/vm"
)

func Import(env *vm.Env) *vm.Env {
	panic("can't import 'github.com/daviddengcn/go-colortext'")
	return nil
}
