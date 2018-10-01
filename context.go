package buffalo

import (
	"context"
	"net/http"
	"sync"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/x/httpx"
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

// GetRouteHelpers returns a map of BuildPathHelper() for each route available in the app.
func (a *App) GetRouteHelpers() map[string]interface{} {
	routeHelpers := map[string]interface{}{}

	for _, route := range a.Routes() {
		cRoute := route
		routeHelpers[cRoute.PathName] = cRoute.BuildPathHelper()
	}
	return routeHelpers
}

func (a *App) newContext(info RouteInfo, res http.ResponseWriter, req *http.Request) Context {
	if ws, ok := res.(*Response); ok {
		res = ws
	}
	params := req.URL.Query()
	vars := mux.Vars(req)
	for k, v := range vars {
		params.Set(k, v)
	}

	if err := req.ParseForm(); err == nil {
		for k, v := range req.Form {
			for _, vv := range v {
				params.Set(k, vv)
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

	routeHelpers := a.GetRouteHelpers()
	for helper, route := range routeHelpers {
		contextData[helper] = route
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
