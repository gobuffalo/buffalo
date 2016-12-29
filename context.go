package buffalo

import (
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/gobuffalo/buffalo/render"
)

// Context holds on to information as you
// pass it down through middleware, Handlers,
// templates, etc... It strives to make your
// life a happier one.
type Context interface {
	Response() http.ResponseWriter
	Request() *http.Request
	Session() *Session
	Params() ParamValues
	Param(string) string
	ParamInt(string) (int, error)
	Set(string, interface{})
	Get(string) interface{}
	LogField(string, interface{})
	LogFields(map[string]interface{})
	Logger() Logger
	Bind(interface{}) error
	Render(int, render.Renderer) error
	Error(int, error) error
	Websocket() (*websocket.Conn, error)
	Redirect(int, string, ...interface{}) error
	Data() map[string]interface{}
}

// ParamValues will most commonly be url.Values,
// but isn't it great that you set your own? :)
type ParamValues interface {
	Get(string) string
}
