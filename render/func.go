package render

import "io"

// RendererFunc is the interface for the the function
// needed by the Func renderer.
type RendererFunc func(io.Writer, Data) error

type funcRenderer struct {
	contentType string
	renderFunc  RendererFunc
}

// ContentType returns the content type for this render.
// Examples would be "text/html" or "application/json".
func (s funcRenderer) ContentType() string {
	return s.contentType
}

// Render the provided Data to the provider Writer using the
// RendererFunc provide.
func (s funcRenderer) Render(w io.Writer, data Data) error {
	return s.renderFunc(w, data)
}

// Func renderer allows for easily building one of renderers
// using just a RendererFunc and not having to build a whole
// implementation of the Render interface.
func Func(s string, fn RendererFunc) Renderer {
	return funcRenderer{
		contentType: s,
		renderFunc:  fn,
	}
}

// Func renderer allows for easily building one of renderers
// using just a RendererFunc and not having to build a whole
// implementation of the Render interface.
func (e *Engine) Func(s string, fn RendererFunc) Renderer {
	return Func(s, fn)
}
