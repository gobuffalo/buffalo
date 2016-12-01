package buffalo

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

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
