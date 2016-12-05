package render_test

import (
	"bytes"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	r := require.New(t)

	type ji func(v string) render.Renderer

	table := []ji{
		render.String,
		render.New(render.Options{}).String,
	}

	for _, j := range table {
		re := j("{{name}}")
		r.Equal("text/plain", re.ContentType())
		bb := &bytes.Buffer{}
		err := re.Render(bb, map[string]interface{}{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", bb.String())
	}
}
