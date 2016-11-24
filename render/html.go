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
	names []string
}

func (s htmlRenderer) ContentType() string {
	return "text/html"
}

func (s htmlRenderer) Render(w io.Writer, data interface{}) error {
	s.moot.Lock()
	defer s.moot.Unlock()

	names := s.names
	if s.HTMLLayout != "" {
		names = append(names, s.HTMLLayout)
	}

	if len(names) > 1 {
		lname := names[1]
		layout := s.templates.Lookup(lname)
		if layout == nil {
			b, err := ioutil.ReadFile(lname)
			if err != nil {
				return err
			}
			layout, err = s.templates.New(names[0]).Parse(string(b))
			if err != nil {
				return err
			}
		}
		layout = layout.Funcs(template.FuncMap{
			"yield": s.yield(names[0], data),
		})
		return layout.Execute(w, data)
	}

	return s.executeTemplate(names[0], w, data)
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

func HTML(names ...string) Renderer {
	e := New(&Options{})
	return e.HTML(names...)
}

func (e *Engine) HTML(names ...string) Renderer {
	hr := htmlRenderer{
		Engine: e,
		names:  names,
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
