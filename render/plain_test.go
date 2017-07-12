package render_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/plush"
	"github.com/stretchr/testify/require"
)

func Test_Plain(t *testing.T) {
	r := require.New(t)

	tmpFile, err := ioutil.TempFile("", "test")
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("<%= name %>"))
	r.NoError(err)

	type ji func(...string) render.Renderer

	j := render.New(render.Options{
		TemplateEngine: plush.BuffaloRenderer,
	}).Plain

	re := j(tmpFile.Name())
	r.Equal("text/plain; charset=utf-8", re.ContentType())
	var examples = []string{"Mark", "Jém"}
	for _, example := range examples {
		example := example
		bb := &bytes.Buffer{}
		err := re.Render(bb, map[string]interface{}{"name": example})
		r.NoError(err)
		r.Equal(example, bb.String())
	}
}
