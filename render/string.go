package render

import (
	"io"

	"github.com/aymerick/raymond"
)

type stringRenderer struct {
	body string
}

func (s stringRenderer) ContentType() string {
	return "text/plain"
}

func (s stringRenderer) Render(w io.Writer, data Data) error {
	b, err := raymond.Render(s.body, data)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(b))
	return err
}

// String renderer that will run the string through
// the github.com/aymerick/raymond package and return
// "text/plain" as the content type.
func String(s string) Renderer {
	return stringRenderer{body: s}
}

// String renderer that will run the string through
// the github.com/aymerick/raymond package and return
// "text/plain" as the content type.
func (e *Engine) String(s string) Renderer {
	return String(s)
}
