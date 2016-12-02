package buffalo

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler is the basis for all of Buffalo. A Handler
// will be given a Context interface that represents the
// give request/response. It is the responsibility of the
// Handler to handle the request/response correctly. This
// could mean rendering a template, JSON, etc... or it could
// mean returning an error.
/*
	func (c Context) error {
		return C.Render(200, render.String("Hello World!"))
	}

	func (c Context) error {
		return C.Redirect(301, "http://github.com/markbates/buffalo")
	}

	func (c Context) error {
		return c.Error(422, errors.New("oops!!"))
	}
*/
type Handler func(Context) error

func (a *App) handlerToHandler(h Handler) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ws := res.(*buffaloResponse)
		params := req.URL.Query()
		for _, v := range p {
			params.Set(v.Key, v.Value)
		}

		c := &DefaultContext{
			response: ws,
			request:  req,
			params:   params,
			logger:   a.Logger,
			session:  a.getSession(req, ws),
			data:     map[string]interface{}{},
		}

		err := a.middlewareStack.handler(h)(c)

		if err != nil {
			c.Error(500, err)
		}
	}
}
