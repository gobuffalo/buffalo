package middleware

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
)

// SessionSaver is deprecated, and will be removed in v0.10.0. This now happens automatically. This middleware is no longer required.
func SessionSaver(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		warningMsg := "SessionSaver is deprecated, and will be removed in v0.10.0. This now happens automatically. This middleware is no longer required."
		fmt.Println(warningMsg)
		return next(c)
	}
}
