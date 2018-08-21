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
	return WrapHandler(h)
}

// WrapBuffaloHandler wraps a buffalo.Handler to
// standard http.Handler
func WrapBuffaloHandler(h Handler) http.Handler {
	a := New(Options{})
	// it doesn't matter what we actually map it
	// GET, POST, etc... we just need the underlying
	// RouteInfo, which implements http.Handler
	ri := a.GET("/", h)
	return ri
}

// WrapBuffaloHandlerFunc wraps a buffalo.Handler to
// standard http.HandlerFunc
func WrapBuffaloHandlerFunc(h Handler) http.HandlerFunc {
	return WrapBuffaloHandler(h).ServeHTTP
}
