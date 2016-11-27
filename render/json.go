package render

import (
	"encoding/json"
	"io"
)

type jsonRenderer struct {
	value interface{}
}

func (s jsonRenderer) ContentType() string {
	return "application/json"
}

func (s jsonRenderer) Render(w io.Writer, data Data) error {
	return json.NewEncoder(w).Encode(s.value)
}

func JSON(v interface{}) Renderer {
	return jsonRenderer{value: v}
}

func (e *Engine) JSON(v interface{}) Renderer {
	return JSON(v)
}
