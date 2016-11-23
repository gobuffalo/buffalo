package render

import (
	"encoding/xml"
	"io"
)

type xmlRenderer struct {
	value interface{}
}

func (s xmlRenderer) ContentType() string {
	return "application/xml"
}

func (s xmlRenderer) Render(w io.Writer, data interface{}) error {
	return xml.NewEncoder(w).Encode(s.value)
}

func XML(v interface{}) Renderer {
	return xmlRenderer{value: v}
}

func (e *Engine) XML(v interface{}) Renderer {
	return XML(v)
}
