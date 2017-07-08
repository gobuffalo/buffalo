package render_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
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
	r.Equal("text/plain", re.ContentType())
	bb := &bytes.Buffer{}
	err = re.Render(bb, map[string]interface{}{"name": "Mark"})
	r.NoError(err)
	r.Equal("Mark", strings.TrimSpace(bb.String()))
}
