package buffalo

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Handler is the basis for all of Buffalo. A Handler
// will be given a Context interface that represents the
// give request/response. It is the responsibility of the
// Handler to handle the request/response correctly. This
// could mean rendering a template, JSON, etc... or it could
// mean returning an error.
/*
	func (c Context) error {
		return c.Render(200, render.String("Hello World!"))
	}

	func (c Context) error {
		return c.Redirect(301, "http://github.com/gobuffalo/buffalo")
	}

	func (c Context) error {
		return c.Error(422, errors.New("oops!!"))
	}
*/
type Handler func(Context) error

func (a *App) newContext(info RouteInfo, res http.ResponseWriter, req *http.Request) Context {
	ws := res.(*Response)
	params := req.URL.Query()
	vars := mux.Vars(req)
	for k, v := range vars {
		params.Set(k, v)
	}

	session := a.getSession(req, ws)

	contextData := map[string]interface{}{
		"env":           a.Env,
		"routes":        a.Routes(),
		"current_route": info,
		"current_path":  req.URL.Path,
	}

	for _, route := range a.Routes() {
		cRoute := route
		contextData[cRoute.PathName] = cRoute.BuildPathHelper()
	}

	return &DefaultContext{
		Context:  req.Context(),
		response: ws,
		request:  req,
		params:   params,
		logger:   a.Logger,
		session:  session,
		flash:    newFlash(session),
		data:     contextData,
	}
}

func (info RouteInfo) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if ok := info.processPreHandlers(res, req); !ok {
		return
	}

	a := info.App

	c := a.newContext(info, res, req)

	defer c.Flash().persist(c.Session())

	err := a.Middleware.handler(info)(c)

	if err != nil {
		status := 500
		// unpack root cause and check for HTTPError
		cause := errors.Cause(err)
		httpError, ok := cause.(HTTPError)
		if ok {
			status = httpError.Status
		}
		eh := a.ErrorHandlers.Get(status)
		err = eh(status, err, c)
		if err != nil {
			// things have really hit the fan if we're here!!
			a.Logger.Error(err)
			c.Response().WriteHeader(500)
			c.Response().Write([]byte(err.Error()))
		}
	}
}

func (info RouteInfo) processPreHandlers(res http.ResponseWriter, req *http.Request) bool {
	a := info.App

	sh := func(h http.Handler) bool {
		h.ServeHTTP(res, req)
		if br, ok := res.(*Response); ok {
			if (br.Status < 200 || br.Status > 299) && br.Status > 0 {
				return false
			}
			if br.Size > 0 {
				return false
			}
		}
		return true
	}

	for _, ph := range a.PreHandlers {
		if ok := sh(ph); !ok {
			return false
		}
	}

	last := http.Handler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {}))
	if len(a.PreWares) > 0 {
		for _, ph := range a.PreWares {
			last = ph(last)
			if ok := sh(last); !ok {
				return false
			}
		}
	}
	return true
}
