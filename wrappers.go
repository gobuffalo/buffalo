package buffalo

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/gobuffalo/buffalo/internal/httpx"
	"github.com/gorilla/mux"
)

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

// WrapBuffaloHandler wraps a buffalo.Handler to a standard http.Handler
//
// NOTE: A buffalo Handler expects a buffalo Context. WrapBuffaloHandler uses
// the same logic as DefaultContext where possible, but some functionality
// (e.g. sessions and logging) WILL NOT work with this unwrap function. If
// those features are needed a custom UnwrapHandlerFunc needs to be
// implemented that provides a Context implementing those features.
func WrapBuffaloHandler(h Handler) http.Handler {
	return WrapBuffaloHandlerFunc(h)
}

// WrapBuffaloHandlerFunc wraps a buffalo.Handler to a standard http.HandlerFunc
//
// NOTE: A buffalo Handler expects a buffalo Context. WrapBuffaloHandlerFunc uses
// the same logic as DefaultContext where possible, but some functionality
// (e.g. sessions and logging) WILL NOT work with this unwrap function. If
// those features are needed a custom WrapBuffaloHandlerFunc needs to be
// implemented that provides a Context implementing those features.
func WrapBuffaloHandlerFunc(h Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if ws, ok := res.(*Response); ok {
			res = ws
		}

		// Parse URL Params
		params := url.Values{}
		vars := mux.Vars(req)
		for k, v := range vars {
			params.Add(k, v)
		}

		// Parse URL Query String Params
		// For POST, PUT, and PATCH requests, it also parse the request body as a form.
		// Request body parameters take precedence over URL query string values in params
		if err := req.ParseForm(); err == nil {
			for k, v := range req.Form {
				for _, vv := range v {
					params.Add(k, vv)
				}
			}
		}

		ct := httpx.ContentType(req)

		data := &sync.Map{}
		data.Store("current_path", req.URL.Path)
		data.Store("contentType", ct)
		data.Store("method", req.Method)

		c := &DefaultContext{
			Context:     req.Context(),
			contentType: ct,
			response:    res,
			request:     req,
			params:      params,
			flash:       &Flash{data: map[string][]string{}},
			data:        data,
		}
		h(c)
	}
}
