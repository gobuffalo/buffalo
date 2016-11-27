package render

import "io"

type RenderFunc func(io.Writer, Data) error

type funcRenderer struct {
	contentType string
	renderFunc  RenderFunc
}

func (s funcRenderer) ContentType() string {
	return s.contentType
}

func (s funcRenderer) Render(w io.Writer, data Data) error {
	return s.renderFunc(w, data)
}

func Func(s string, fn RenderFunc) Renderer {
	return funcRenderer{
		contentType: s,
		renderFunc:  fn,
	}
}

func (e *Engine) Func(s string, fn RenderFunc) Renderer {
	return Func(s, fn)
}
