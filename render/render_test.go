package render

import (
	"os"
	"testing"
	"testing/fstest"
)

type Widget struct {
	Name string
}

func (w Widget) ToPath() string {
	return w.Name
}

func NewEngine() *Engine {
	return New(Options{
		TemplatesFS: fstest.MapFS{},
		AssetsFS:    fstest.MapFS{},
	})
}

type rendFriend func(string, RendererFunc) Renderer

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func init() {
	assetMap.Range(func(key, value string) bool {
		assetMap.Delete(key)
		return true
	})
}
