package render

import (
	"fmt"
	"io"

	"errors"
)

type stringRenderer struct {
	*Engine
	body string
}

func (s stringRenderer) ContentType() string {
	return "text/plain; charset=utf-8"
}

func (s stringRenderer) Render(w io.Writer, data Data) error {
	te, ok := s.TemplateEngines["text"]
	if !ok {
		return errors.New("could not find a template engine for text")
	}
	t, err := te(s.body, data, s.Helpers)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(t))
	return err
}

// String renderer that will run the string through
// the github.com/gobuffalo/plush package and return
// "text/plain" as the content type.
func String(s string, args ...interface{}) Renderer {
	e := New(Options{})
	return e.String(s, args...)
}

// String renderer that will run the string through
// the github.com/gobuffalo/plush package and return
// "text/plain" as the content type.
func (e *Engine) String(s string, args ...interface{}) Renderer {
	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}
	return stringRenderer{
		Engine: e,
		body:   s,
	}
}
