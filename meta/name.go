package meta

import (
	"fmt"
	"runtime"

	"github.com/markbates/inflect"
)

// Name is deprecated, please use github.com/markbates/inflect.Name instead.
func Name(s string) inflect.Name {
	warningMsg := "Name is deprecated, and will be removed in v0.12.0. Use github.com/markbates/inflect.Name instead."
	if _, file, no, ok := runtime.Caller(1); ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}
	fmt.Println(warningMsg)
	return inflect.Name(s)
}
