// +build !appengine

package os

import (
	"github.com/mattn/anko/vm"
	pkg "os"
	"reflect"
)

func handleAppEngine(m *vm.Env) {
	m.Define("Getppid", reflect.ValueOf(pkg.Getppid))
}
