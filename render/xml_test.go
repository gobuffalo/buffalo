package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_XML(t *testing.T) {
	r := require.New(t)

	type ji func(v interface{}) render.Renderer

	table := []ji{
		render.XML,
		render.New(render.Options{}).XML,
	}

	type user struct {
		Name string
	}

	for _, j := range table {
		re := j(user{Name: "mark"})
		r.Equal("application/xml", re.ContentType())
		bb := &bytes.Buffer{}
		err := re.Render(bb, nil)
		r.NoError(err)
		r.Equal(`<user><Name>mark</Name></user>`, strings.TrimSpace(bb.String()))
	}
}
