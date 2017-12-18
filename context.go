package buffalo

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/buffalo/render"
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
	Websocket() (*websocket.Conn, error)
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
