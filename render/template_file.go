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

func (s *templateFileRenderer) Render(w io.Writer, data Data) error {
	s.moot.Lock()
	defer s.moot.Unlock()

	names := s.names

	tname := names[0]
	tm, err := s.Lookup(tname)
	if err != nil {
		return errors.WithStack(err)
	}

	if len(names) > 1 {
		tname = names[1]
		tm, err = s.Lookup(tname)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	tm = tm.Funcs(template.FuncMap{
		"yield":   s.yield(names[0], data),
		"partial": s.partial(data),
	})

	return tm.Execute(w, data)
}

func TemplateFile(c string, names ...string) Renderer {
	e := New(&Options{})
	return e.TemplateFile(c, names...)
}

func (e *Engine) TemplateFile(c string, names ...string) Renderer {
	return &templateFileRenderer{
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
			return s.htmlError(err)
		}
		return template.HTML(bb.String())
	}
}

func (s *templateFileRenderer) partial(data Data) func(string) template.HTML {
	return func(name string) template.HTML {
		d, f := filepath.Split(name)
		name = filepath.Join(d, "_"+f)
		return s.yield(name, data)()
	}
}

func (s templateFileRenderer) executeTemplate(name string, w io.Writer, data Data) error {
	tm, err := s.Lookup(name)
	if err != nil {
		return err
	}
	return tm.Execute(w, data)
}

func (s templateFileRenderer) Lookup(name string) (*template.Template, error) {
	tm := s.templates.Lookup(name)
	if tm == nil {
		b, err := ioutil.ReadFile(filepath.Join(s.TemplatesPath, name))
		if err != nil {
			return tm, errors.WithStack(fmt.Errorf("could not find template: %s", name))
		}
		tm, err = template.New(name).Funcs(s.TemplateFuncs).Parse(string(b))
		if err != nil {
			return tm, errors.WithStack(err)
		}
	}
	return tm, nil
}

func (s templateFileRenderer) htmlError(err error) template.HTML {
	return template.HTML(fmt.Sprintf("<pre>%s</pre>", errors.WithStack(err).Error()))
}
