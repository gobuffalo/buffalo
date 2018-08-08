package buffalo

import (
	"database/sql"
	"net/http"

	gcontext "github.com/gorilla/context"
	"github.com/pkg/errors"
)

func (info RouteInfo) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer gcontext.Clear(req)
	a := info.App

	c := a.newContext(info, res, req)

	defer c.Flash().persist(c.Session())

	err := a.Middleware.handler(info)(c)

	if err != nil {
		status := 500
		// unpack root cause and check for HTTPError
		cause := errors.Cause(err)
		switch cause {
		case sql.ErrNoRows:
			status = 404
		default:
			if h, ok := cause.(HTTPError); ok {
				status = h.Status
			}
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
