package render

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/pkg/errors"
)

type templateRenderer struct {
	*Engine
	contentType string
	names       []string
}

func (s templateRenderer) ContentType() string {
	return s.contentType
}

func (s *templateRenderer) Render(w io.Writer, data Data) error {
	var yield raymond.SafeString
	var err error
	for _, name := range s.names {
		yield, err = s.execute(name, data)
		if err != nil {
			return errors.WithMessage(errors.WithStack(err), name)
		}
		data["yield"] = yield
	}
	_, err = w.Write([]byte(yield))
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *templateRenderer) execute(name string, data Data) (raymond.SafeString, error) {
	source, err := s.source(name)
	if err != nil {
		return raymond.SafeString(fmt.Sprintf("<pre>%s: %s</pre>", name, err.Error())), err
	}
	source.RegisterHelper("partial", func(name string, options *raymond.Options) raymond.SafeString {
		d := data
		for k, v := range options.Hash() {
			d[k] = v
			defer delete(data, k)
		}
		p, err := s.partial(name, d)
		if err != nil {
			return raymond.SafeString(fmt.Sprintf("<pre>%s: %s</pre>", name, err.Error()))
		}
		return p
	})
	yield, err := source.Exec(data)
	if err != nil {
		return raymond.SafeString(fmt.Sprintf("<pre>%s: %s</pre>", name, err.Error())), err
	}
	return raymond.SafeString(yield), nil
}

func (s *templateRenderer) source(name string) (*raymond.Template, error) {
	b, err := ioutil.ReadFile(filepath.Join(s.TemplatesPath, name))
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("could not find template: %s", name))
	}
	t, err := raymond.Parse(string(b))
	if err != nil {
		return t, errors.WithMessage(errors.WithStack(err), name)
	}
	t.RegisterHelpers(s.Helpers)
	return t, err
}

func (s *templateRenderer) partial(name string, data Data) (raymond.SafeString, error) {
	d, f := filepath.Split(name)
	name = filepath.Join(d, "_"+f)
	return s.execute(name, data)
}

// Template renders the named files using the specified
// content type and the github.com/aymerick/raymond
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "{{yield}}".
func Template(c string, names ...string) Renderer {
	e := New(Options{})
	return e.Template(c, names...)
}

// Template renders the named files using the specified
// content type and the github.com/aymerick/raymond
// package for templating. If more than 1 file is provided
// the second file will be considered a "layout" file
// and the first file will be the "content" file which will
// be placed into the "layout" using "{{yield}}".
func (e *Engine) Template(c string, names ...string) Renderer {
	return &templateRenderer{
		Engine:      e,
		contentType: c,
		names:       names,
	}
}
