package render

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

type htmlRenderer struct {
	*Engine
	name string
}

func (s htmlRenderer) ContentType() string {
	return "text/html"
}

func (s htmlRenderer) Render(w io.Writer, data interface{}) error {
	s.moot.Lock()
	defer s.moot.Unlock()

	if s.name != s.HTMLLayout && s.HTMLLayout != "" {
		layout := s.templates.Lookup(s.HTMLLayout)
		if layout == nil {
			b, err := ioutil.ReadFile(s.HTMLLayout)
			if err != nil {
				return err
			}
			layout, err = s.templates.New(s.name).Parse(string(b))
			if err != nil {
				return err
			}
		}
		layout = layout.Funcs(template.FuncMap{
			"yield": s.yield(s.name, data),
		})
		return layout.Execute(w, data)
	}

	return s.executeTemplate(s.name, w, data)
}

func (s htmlRenderer) executeTemplate(name string, w io.Writer, data interface{}) error {
	var t *template.Template

	if t == nil {
		b, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}
		t, err = s.templates.New(name).Parse(string(b))
		if err != nil {
			return err
		}
	}
	s.templates = t
	return t.Execute(w, data)
}

func HTML(name string) Renderer {
	e := New(&Options{})
	return e.HTML(name)
}

func (e *Engine) HTML(name string) Renderer {
	hr := htmlRenderer{
		Engine: e,
		name:   name,
	}
	return hr
}

func (s htmlRenderer) yield(name string, data interface{}) func() template.HTML {
	return func() template.HTML {
		bb := &bytes.Buffer{}
		err := s.executeTemplate(name, bb, data)
		if err != nil {
			return template.HTML(errors.WithStack(err).Error())
		}
		return template.HTML(bb.String())
	}
}
