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

	for _, n := range names {
		t := s.templates.Lookup(n)
		if t == nil {
			b, err := ioutil.ReadFile(n)
			if err != nil {
				return err
			}
			s.templates, err = s.templates.New(n).Parse(string(b))
			if err != nil {
				return err
			}
		}
	}

	if len(names) > 1 {
		lname := names[1]
		layout := s.templates.Lookup(lname)
		layout = layout.Funcs(template.FuncMap{
			"yield": s.yield(names[0], data),
		})
		return layout.Execute(w, data)
	}

	return s.executeTemplate(names[0], w, data)
}

func (s htmlRenderer) executeTemplate(name string, w io.Writer, data interface{}) error {
	return s.templates.ExecuteTemplate(w, name, data)
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
