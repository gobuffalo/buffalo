package render_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_TemplateFile(t *testing.T) {
	r := require.New(t)

	tmpFile, err := ioutil.TempFile("", "test")
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("{{.}}"))
	r.NoError(err)

	type ji func(string, string) render.Renderer

	table := []ji{
		render.TemplateFile,
		render.New(&render.Options{}).TemplateFile,
	}

	for _, j := range table {
		re := j("foo/bar", tmpFile.Name())
		r.Equal("foo/bar", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, "Mark")
		r.NoError(err)
		r.Equal("Mark", strings.TrimSpace(bb.String()))
	}
}
