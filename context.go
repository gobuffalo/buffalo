package buffalo

import (
	"net/http"

	"github.com/markbates/buffalo/render"
)

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
	NoContent(int) error
}

type ParamValues interface {
	Get(string) string
}
