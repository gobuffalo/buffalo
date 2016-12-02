package buffalo

import "net/http"

// WrapHandler wraps a standard http.Handler and transforms it
// into a buffalo.Handler.
func WrapHandler(h http.Handler) Handler {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// WrapHandlerFunc wraps a standard http.HandlerFunc and
// transforms it into a buffalo.Handler.
func WrapHandlerFunc(h http.HandlerFunc) Handler {
	return WrapHandler(http.HandlerFunc(h))
}
