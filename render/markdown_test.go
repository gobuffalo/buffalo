package render_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_Markdown(t *testing.T) {
	r := require.New(t)

	tmpDir := filepath.Join(os.TempDir(), "markdown_test")
	err := os.MkdirAll(tmpDir, 0766)
	r.NoError(err)
	defer os.Remove(tmpDir)

	tmpFile, err := os.Create(filepath.Join(tmpDir, "t.md"))
	r.NoError(err)

	_, err = tmpFile.Write([]byte("<%= name %>"))
	r.NoError(err)

	type ji func(...string) render.Renderer
	t.Run("without a layout", func(st *testing.T) {
		r := require.New(st)

		table := []ji{
			render.New(render.Options{
				TemplatesBox: packr.NewBox(tmpDir),
			}).HTML,
		}

		for _, j := range table {
			re := j(filepath.Base(tmpFile.Name()))
			r.Equal("text/html; charset=utf-8", re.ContentType())
			bb := &bytes.Buffer{}
			err = re.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("<p>Mark</p>", strings.TrimSpace(bb.String()))
		}
	})

	t.Run("with a layout", func(st *testing.T) {
		r := require.New(st)

		layout, err := os.Create(filepath.Join(tmpDir, "test.html"))
		r.NoError(err)
		defer os.Remove(layout.Name())

		_, err = layout.Write([]byte("<body><%= yield %></body>"))
		r.NoError(err)

		re := render.New(render.Options{
			HTMLLayout:   filepath.Base(layout.Name()),
			TemplatesBox: packr.NewBox(tmpDir),
		}).HTML(filepath.Base(tmpFile.Name()))

		r.Equal("text/html; charset=utf-8", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, map[string]interface{}{"name": "Mark"})
		r.NoError(err)
		r.Equal("<body><p>Mark</p>\n</body>", strings.TrimSpace(bb.String()))
	})
}
