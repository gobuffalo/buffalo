package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

type templateFileRenderer struct {
	*Engine
	contentType string
	names       []string
}

func (s templateFileRenderer) ContentType() string {
	return s.contentType
}

func (s templateFileRenderer) Render(w io.Writer, data Data) error {
	s.moot.Lock()
	defer s.moot.Unlock()

	names := s.names
	for _, n := range names {
		t := s.templates.Lookup(n)
		if t == nil {
			b, err := ioutil.ReadFile(filepath.Join(s.TemplatesPath, n))
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

func TemplateFile(c string, names ...string) Renderer {
	e := New(&Options{})
	return e.TemplateFile(c, names...)
}

func (e *Engine) TemplateFile(c string, names ...string) Renderer {
	return templateFileRenderer{
		Engine:      e,
		contentType: c,
		names:       names,
	}
}

func (s templateFileRenderer) yield(name string, data Data) func() template.HTML {
	return func() template.HTML {
		bb := &bytes.Buffer{}
		err := s.executeTemplate(name, bb, data)
		if err != nil {
			return template.HTML(fmt.Sprintf("<pre>%s</pre>", errors.WithStack(err).Error()))
		}
		return template.HTML(bb.String())
	}
}

func (s templateFileRenderer) executeTemplate(name string, w io.Writer, data Data) error {
	return s.templates.ExecuteTemplate(w, name, data)
}
