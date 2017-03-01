package render

import "io"

type stringRenderer struct {
	*Engine
	body string
}

func (s stringRenderer) ContentType() string {
	return "text/plain"
}

func (s stringRenderer) Render(w io.Writer, data Data) error {
	t, err := s.TemplateEngine(s.body, data, s.Helpers)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(t))
	return err
}

// String renderer that will run the string through
// the github.com/aymerick/raymond package and return
// "text/plain" as the content type.
func String(s string) Renderer {
	e := New(Options{})
	return e.String(s)
}

// String renderer that will run the string through
// the github.com/aymerick/raymond package and return
// "text/plain" as the content type.
func (e *Engine) String(s string) Renderer {
	return stringRenderer{
		Engine: e,
		body:   s,
	}
}
