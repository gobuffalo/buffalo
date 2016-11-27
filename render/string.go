package render

import (
	"html/template"
	"io"
)

type stringRenderer struct {
	body string
}

func (s stringRenderer) ContentType() string {
	return "text/plain"
}

func (s stringRenderer) Render(w io.Writer, data Data) error {
	t, err := template.New("").Parse(s.body)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

func String(s string) Renderer {
	return stringRenderer{body: s}
}

func (e *Engine) String(s string) Renderer {
	return String(s)
}
