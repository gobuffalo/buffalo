package buffalo

import (
	"net/http"

	"github.com/gobuffalo/buffalo/events"
	gcontext "github.com/gorilla/context"
)

func (info RouteInfo) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	defer gcontext.Clear(req)

	a := info.App

	c := a.newContext(info, res, req)
	defer c.Flash().persist(c.Session())

	payload := map[string]interface{}{
		"route":   info,
		"req":     req,
		"context": c,
	}

	events.Emit(events.Event{
		Kind:    events.RouteStarted,
		Payload: payload,
	})

	err := a.Middleware.handler(info)(c)

	if err != nil {
		events.Emit(events.Event{
			Kind:    events.ErrRoute,
			Payload: payload,
			Error:   err,
		})
		// things have really hit the fan if we're here!!
		a.Logger.Error(err)
		c.Response().WriteHeader(500)
		c.Response().Write([]byte(err.Error()))
	}
	events.Emit(events.Event{
		Kind:    events.RouteFinished,
		Payload: payload,
	})
}
