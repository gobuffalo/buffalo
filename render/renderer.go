package render

import "io"

// Renderer interface that must be satisified to be used with
// buffalo.Context.Render
type Renderer interface {
	ContentType() string
	Render(io.Writer, Data) error
}

// Data type to be provided to the Render function on the
// Renderer interface.
type Data map[string]interface{}
