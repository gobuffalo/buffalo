package render

import (
	"encoding/xml"
	"io"
)

type xmlRenderer struct {
	value interface{}
}

func (s xmlRenderer) ContentType() string {
	return "application/xml; charset=utf-8"
}

func (s xmlRenderer) Render(w io.Writer, data Data) error {
	io.WriteString(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	return enc.Encode(s.value)
}

// XML renders the value using the "application/xml"
// content type.
func XML(v interface{}) Renderer {
	return xmlRenderer{value: v}
}

// XML renders the value using the "application/xml"
// content type.
func (e *Engine) XML(v interface{}) Renderer {
	return XML(v)
}
