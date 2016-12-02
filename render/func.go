package render

import "io"

type RendererFunc func(io.Writer, Data) error

type funcRenderer struct {
	contentType string
	renderFunc  RendererFunc
}

func (s funcRenderer) ContentType() string {
	return s.contentType
}

func (s funcRenderer) Render(w io.Writer, data Data) error {
	return s.renderFunc(w, data)
}

func Func(s string, fn RendererFunc) Renderer {
	return funcRenderer{
		contentType: s,
		renderFunc:  fn,
	}
}

func (e *Engine) Func(s string, fn RendererFunc) Renderer {
	return Func(s, fn)
}
