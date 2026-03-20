package render_test

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	r := require.New(t)

	j := render.New(render.Options{}).String

	re := j("<%= name %>")
	r.Equal("text/plain; charset=utf-8", re.ContentType())

	var examples = []string{"Mark", "Jém"}
	for _, example := range examples {
		bb := &bytes.Buffer{}
		err := re.Render(bb, map[string]any{"name": example})
		r.NoError(err)
		r.Equal(example, bb.String())
	}

}
