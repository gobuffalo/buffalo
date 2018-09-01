package buffalo

import (
	"net/http"

	gcontext "github.com/gorilla/context"
)

func (info RouteInfo) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer gcontext.Clear(req)
	a := info.App

	c := a.newContext(info, res, req)

	defer c.Flash().persist(c.Session())

	err := a.Middleware.handler(info)(c)

	if err != nil {
		// things have really hit the fan if we're here!!
		a.Logger.Error(err)
		c.Response().WriteHeader(500)
		c.Response().Write([]byte(err.Error()))
	}
}
