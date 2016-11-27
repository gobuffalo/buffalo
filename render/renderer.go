package render

import "io"

type Renderer interface {
	ContentType() string
	Render(io.Writer, Data) error
}

type Data map[string]interface{}
