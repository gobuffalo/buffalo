package render_test

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_Template(t *testing.T) {
	r := require.New(t)

	type ji func(string, *template.Template) render.Renderer

	table := []ji{
		render.Template,
		render.New(render.Options{}).Template,
	}

	for _, j := range table {
		tmp, err := template.New("").Parse("{{.name}}")
		r.NoError(err)
		re := j("foo/bar", tmp)
		r.Equal("foo/bar", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, map[string]interface{}{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", strings.TrimSpace(bb.String()))
	}
}
