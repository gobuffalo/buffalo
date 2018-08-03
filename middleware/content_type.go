package middleware

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
)

// SetContentType on the request to desired type. This will
// override any content type sent by the client.
//
// Deprecated: use github.com/gobuffalo/mw-contenttype#Set instead.
func SetContentType(s string) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			fmt.Printf("SetContentType middleware is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-contenttype#Set instead.")
			c.Request().Header.Set("Content-Type", s)
			return next(c)
		}
	}
}

// AddContentType will add a secondary content type to
// a request. If no content type is sent by the client
// the default will be set, otherwise the client's
// content type will be used.
//
// Deprecated: use github.com/gobuffalo/mw-contenttype#Add instead.
func AddContentType(s string) buffalo.MiddlewareFunc {
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			fmt.Printf("SetContentType middleware is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-contenttype#Add instead.")
			c.Request().Header.Add("Content-Type", s)
			return next(c)
		}
	}
}
