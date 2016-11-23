package render

import (
	"html/template"
	"io"
	"io/ioutil"
)

type templateFileRenderer struct {
	contentType string
	path        string
}

func (s templateFileRenderer) ContentType() string {
	return s.contentType
}

func (s templateFileRenderer) Render(w io.Writer, data interface{}) error {
	b, err := ioutil.ReadFile(s.path)
	if err != nil {
		return err
	}
	t, err := template.New(s.path).Parse(string(b))
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

func TemplateFile(c string, path string) Renderer {
	return templateFileRenderer{contentType: c, path: path}
}

func (e *Engine) TemplateFile(c string, path string) Renderer {
	return TemplateFile(c, path)
}
