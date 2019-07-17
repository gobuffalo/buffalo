package render

import (
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
		AssetsBox:    packd.NewMemoryBox(),
	})
}

type rendFriend func(string, RendererFunc) Renderer
