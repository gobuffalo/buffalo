package buffalo

import (
	"fmt"
	"log"
	"runtime"

	"github.com/gobuffalo/buffalo/binding"
)

// RegisterBinder is deprecated. Please use binding.Register instead.
func RegisterBinder(contentType string, fn binding.Binder) {
	warningMsg := "RegisterBinder is deprecated, and will be removed in v0.10.0. Use binding.Register instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}

	log.Println(warningMsg)

	binding.Register(contentType, fn)
}
