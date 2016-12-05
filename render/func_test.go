package render_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_Func(t *testing.T) {
	r := require.New(t)

	type ji func(string, render.RendererFunc) render.Renderer

	table := []ji{
		render.Func,
		render.New(render.Options{}).Func,
	}

	for _, j := range table {
		bb := &bytes.Buffer{}
		re := j("foo/bar", func(w io.Writer, data render.Data) error {
			_, err := w.Write([]byte(data["name"].(string)))
			return err
		})

		r.Equal("foo/bar", re.ContentType())
		err := re.Render(bb, render.Data{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", bb.String())
	}
}
