package render

import (
	"github.com/psanford/memfs"
)

type Widget struct {
	Name string
}

func (w Widget) ToPath() string {
	return w.Name
}

func NewEngine() *Engine {
	return New(Options{
		TemplatesFS: memfs.New(),
		AssetsFS:    memfs.New(),
	})
}

type rendFriend func(string, RendererFunc) Renderer
