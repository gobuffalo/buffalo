package buffalo

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/buffalo/internal/httpx"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gorilla/mux"
)

// Context holds on to information as you
// pass it down through middleware, Handlers,
// templates, etc... It strives to make your
// life a happier one.
type Context interface {
	context.Context
	Response() http.ResponseWriter
	Request() *http.Request
	Session() *Session
	Cookies() *Cookies
	Params() ParamValues
	Param(string) string
	Set(string, interface{})
	LogField(string, interface{})
	LogFields(map[string]interface{})
	Logger() Logger
	Bind(interface{}) error
	Render(int, render.Renderer) error
	Error(int, error) error
	Redirect(int, string, ...interface{}) error
	Data() map[string]interface{}
	Flash() *Flash
	File(string) (binding.File, error)
}

// ParamValues will most commonly be url.Values,
// but isn't it great that you set your own? :)
type ParamValues interface {
	Get(string) string
}

func (a *App) newContext(info RouteInfo, res http.ResponseWriter, req *http.Request) Context {
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

	session := a.getSession(req, res)

	ct := httpx.ContentType(req)
	contextData := map[string]interface{}{
		"app":           a,
		"env":           a.Env,
		"routes":        a.Routes(),
		"current_route": info,
		"current_path":  req.URL.Path,
		"contentType":   ct,
		"method":        req.Method,
	}

	for _, route := range a.Routes() {
		cRoute := route
		contextData[cRoute.PathName] = cRoute.BuildPathHelper()
	}

	return &DefaultContext{
		Context:     req.Context(),
		contentType: ct,
		response:    res,
		request:     req,
		params:      params,
		logger:      a.Logger,
		session:     session,
		flash:       newFlash(session),
		data:        contextData,
		moot:        &sync.RWMutex{},
	}
}
