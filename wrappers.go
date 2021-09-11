package buffalo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/buffalo/internal/httpx"
	"github.com/gobuffalo/buffalo/internal/takeon/github.com/markbates/errx"
	"github.com/gobuffalo/buffalo/render"
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

		c := &wrappedContext{
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

type wrappedContext struct {
	context.Context
	response    http.ResponseWriter
	request     *http.Request
	params      url.Values
	contentType string
	flash       *Flash
	data        *sync.Map
}

func (w *wrappedContext) Response() http.ResponseWriter     { return w.response }
func (w *wrappedContext) Request() *http.Request            { return w.request }
func (w *wrappedContext) Session() *Session                 { return nil }
func (w *wrappedContext) Cookies() *Cookies                 { return &Cookies{w.request, w.response} }
func (w *wrappedContext) Params() ParamValues               { return w.params }
func (w *wrappedContext) Param(key string) string           { return w.Params().Get(key) }
func (w *wrappedContext) Set(key string, value interface{}) { w.data.Store(key, value) }
func (w *wrappedContext) LogField(string, interface{})      {}
func (w *wrappedContext) LogFields(map[string]interface{})  {}
func (w *wrappedContext) Logger() Logger                    { return nil }
func (w *wrappedContext) Bind(value interface{}) error      { return binding.Exec(w.Request(), value) }
func (w *wrappedContext) Flash() *Flash                     { return w.flash }
func (w *wrappedContext) Error(status int, err error) error {
	return HTTPError{Status: status, Cause: err}
}

func (w *wrappedContext) Render(status int, rr render.Renderer) error {
	if rr == nil {
		w.Response().WriteHeader(status)
		return nil
	}

	data := w.Data()
	pp := map[string]string{}
	for k, v := range w.params {
		pp[k] = v[0]
	}
	data["params"] = pp
	data["flash"] = w.Flash().data
	data["session"] = w.Session()
	data["request"] = w.Request()
	data["status"] = status
	bb := &bytes.Buffer{}

	err := rr.Render(bb, data)
	if err != nil {
		if er, ok := errx.Unwrap(err).(render.ErrRedirect); ok {
			return w.Redirect(er.Status, er.URL)
		}
		return HTTPError{Status: http.StatusInternalServerError, Cause: err}
	}

	if w.Session() != nil {
		w.Flash().Clear()
		w.Flash().persist(w.Session())
	}

	w.Response().Header().Set("Content-Type", rr.ContentType())
	if p, ok := data["pagination"].(paginable); ok {
		w.Response().Header().Set("X-Pagination", p.Paginate())
	}
	w.Response().WriteHeader(status)
	_, err = io.Copy(w.Response(), bb)
	if err != nil {
		return HTTPError{Status: http.StatusInternalServerError, Cause: err}
	}

	return nil
}

func (w *wrappedContext) Redirect(status int, url string, args ...interface{}) error {
	w.Flash().persist(w.Session())

	if strings.HasSuffix(url, "Path()") {
		if len(args) > 1 {
			return fmt.Errorf("you must pass only a map[string]interface{} to a route path: %T", args)
		}
		var m map[string]interface{}
		if len(args) == 1 {
			rv := reflect.Indirect(reflect.ValueOf(args[0]))
			if !rv.Type().ConvertibleTo(mapType) {
				return fmt.Errorf("you must pass only a map[string]interface{} to a route path: %T", args)
			}
			m = rv.Convert(mapType).Interface().(map[string]interface{})
		}
		h, ok := w.Value(strings.TrimSuffix(url, "()")).(RouteHelperFunc)
		if !ok {
			return fmt.Errorf("could not find a route helper named %s", url)
		}
		url, err := h(m)
		if err != nil {
			return err
		}
		http.Redirect(w.Response(), w.Request(), string(url), status)
		return nil
	}

	if len(args) > 0 {
		url = fmt.Sprintf(url, args...)
	}
	http.Redirect(w.Response(), w.Request(), url, status)
	return nil
}

func (w *wrappedContext) Data() map[string]interface{} {
	m := map[string]interface{}{}
	w.data.Range(func(k, v interface{}) bool {
		s, ok := k.(string)
		if !ok {
			return false
		}
		m[s] = v
		return true
	})
	return m
}

func (w *wrappedContext) File(name string) (binding.File, error) {
	req := w.Request()
	if err := req.ParseMultipartForm(5 * 1024 * 1024); err != nil {
		return binding.File{}, err
	}
	f, h, err := req.FormFile(name)
	bf := binding.File{
		File:       f,
		FileHeader: h,
	}
	return bf, err
}
