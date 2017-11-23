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

func Test_HTML(t *testing.T) {
	r := require.New(t)

	tmpDir := filepath.Join(os.TempDir(), "html_test")
	err := os.MkdirAll(tmpDir, 0766)
	r.NoError(err)
	defer os.Remove(tmpDir)

	tmpFile, err := os.Create(filepath.Join(tmpDir, "test.html"))
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("<%= name %>"))
	r.NoError(err)

	t.Run("without a layout", func(st *testing.T) {
		r := require.New(st)

		j := render.New(render.Options{
			TemplatesBox: packr.NewBox(tmpDir),
		}).HTML

		re := j(filepath.Base(tmpFile.Name()))
		r.Equal("text/html", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, map[string]interface{}{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", strings.TrimSpace(bb.String()))
	})

	t.Run("with a layout", func(st *testing.T) {
		r := require.New(st)

		layout, err := os.Create(filepath.Join(tmpDir, "layout.html"))
		r.NoError(err)
		defer os.Remove(layout.Name())

		_, err = layout.Write([]byte("<body><%= yield %></body>"))
		r.NoError(err)

		re := render.New(render.Options{
			TemplatesBox: packr.NewBox(tmpDir),
			HTMLLayout:   filepath.Base(layout.Name()),
		})

		st.Run("using just the HTMLLayout", func(sst *testing.T) {
			r := require.New(sst)
			h := re.HTML(filepath.Base(tmpFile.Name()))

			r.Equal("text/html", h.ContentType())
			bb := &bytes.Buffer{}
			err = h.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
		})

		st.Run("overriding the HTMLLayout", func(sst *testing.T) {
			r := require.New(sst)
			nlayout, err := os.Create(filepath.Join(tmpDir, "layout2.html"))
			r.NoError(err)
			defer os.Remove(nlayout.Name())

			_, err = nlayout.Write([]byte("<html><%= yield %></html>"))
			r.NoError(err)
			h := re.HTML(filepath.Base(tmpFile.Name()), filepath.Base(nlayout.Name()))

			r.Equal("text/html", h.ContentType())
			bb := &bytes.Buffer{}
			err = h.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("<html>Mark</html>", strings.TrimSpace(bb.String()))
		})

	})
}
