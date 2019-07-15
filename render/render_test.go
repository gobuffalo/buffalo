package render

import (
	"context"
	"net/http"

	"github.com/gobuffalo/packd"
)

type Widget struct {
	Name string
}

func (w Widget) ToPath() string {
	return w.Name
}

func NewEngine() *Engine {
	return New(Options{
		TemplatesBox: packd.NewMemoryBox(),
	})
}

type rendFriend func(string, RendererFunc) Renderer

type testContext struct {
	context.Context
	rw http.ResponseWriter
}

func (c testContext) Response() http.ResponseWriter {
	return c.rw
}
