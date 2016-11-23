package render

import (
	"html/template"
	"io"
)

type templateRenderer struct {
	contentType string
	template    *template.Template
}

func (s templateRenderer) ContentType() string {
	return s.contentType
}

func (s templateRenderer) Render(w io.Writer, data interface{}) error {
	return s.template.Execute(w, data)
}

func Template(c string, t *template.Template) Renderer {
	return templateRenderer{contentType: c, template: t}
}

func (e *Engine) Template(c string, t *template.Template) Renderer {
	return Template(c, t)
}
