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

func withHTMLFile(name string, contents string, fn func(*Engine)) error {
	box := packd.NewMemoryBox()
	box.AddString(name, contents)
	e := New(Options{
		TemplatesBox: box,
	})
	fn(e)
	return nil
}

type rendFriend func(string, RendererFunc) Renderer

type testContext struct {
	context.Context
	rw http.ResponseWriter
}

func (c testContext) Response() http.ResponseWriter {
	return c.rw
}
