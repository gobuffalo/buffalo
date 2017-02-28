package render_test

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/velvet"
	"github.com/stretchr/testify/require"
)

func Test_String(t *testing.T) {
	r := require.New(t)

	j := render.New(render.Options{
		TemplateEngine: velvet.BuffaloRenderer,
	}).String

	re := j("{{name}}")
	r.Equal("text/plain", re.ContentType())
	bb := &bytes.Buffer{}
	err := re.Render(bb, map[string]interface{}{"name": "Mark"})
	r.NoError(err)
	r.Equal("Mark", bb.String())
}
