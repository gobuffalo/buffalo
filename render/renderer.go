package render

import "io"

type Renderer interface {
	ContentType() string
	Render(io.Writer, interface{}) error
}
