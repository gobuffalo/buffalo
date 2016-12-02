package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_JSON(t *testing.T) {
	r := require.New(t)

	type ji func(v interface{}) render.Renderer

	table := []ji{
		render.JSON,
		render.New(render.Options{}).JSON,
	}

	for _, j := range table {
		re := j(map[string]string{"hello": "world"})
		r.Equal("application/json", re.ContentType())
		bb := &bytes.Buffer{}
		err := re.Render(bb, nil)
		r.NoError(err)
		r.Equal(`{"hello":"world"}`, strings.TrimSpace(bb.String()))
	}
}
