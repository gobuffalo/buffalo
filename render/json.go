package render

import (
	"encoding/json"
	"io"
)

type jsonRenderer struct {
	value any
}

func (s jsonRenderer) ContentType() string {
	return "application/json; charset=utf-8"
}

func (s jsonRenderer) Render(w io.Writer, data Data) error {
	return json.NewEncoder(w).Encode(s.value)
}

// JSON renders the value using the "application/json"
// content type.
func JSON(v any) Renderer {
	return jsonRenderer{value: v}
}

// JSON renders the value using the "application/json"
// content type.
func (e *Engine) JSON(v any) Renderer {
	return JSON(v)
}
