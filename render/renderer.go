package render

import (
	"io"

	"github.com/gobuffalo/velvet"
)

// Renderer interface that must be satisfied to be used with
// buffalo.Context.Render
type Renderer interface {
	ContentType() string
	Render(io.Writer, Data) error
}

// Data type to be provided to the Render function on the
// Renderer interface.
type Data map[string]interface{}

// ToVelvet converts the render data into a velvet.Context
func (d Data) ToVelvet() *velvet.Context {
	return velvet.NewContextWith(d)
}
