package render_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_Plain(t *testing.T) {
	r := require.New(t)

	tDir, err := ioutil.TempDir("", "templates")
	if err != nil {
		r.Fail("Could not set the templates dir")
	}

	tmpFile, err := os.Create(filepath.Join(tDir, "test"))
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("<%= name %>"))
	r.NoError(err)

	type ji func(...string) render.Renderer

	j := render.New(render.Options{
		TemplatesBox: packr.NewBox(tDir),
	}).Plain

	re := j(filepath.Base(tmpFile.Name()))
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
