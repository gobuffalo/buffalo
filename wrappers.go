package buffalo

import "net/http"

func WrapHandler(h http.Handler) Handler {
	return func(c Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func WrapHandlerFunc(h http.HandlerFunc) Handler {
	return WrapHandler(http.HandlerFunc(h))
}
