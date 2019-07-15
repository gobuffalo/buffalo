package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HTML(t *testing.T) {
	r := require.New(t)

	err := withHTMLFile("test.html", "<%= name %>", func(e *Engine) {
		t.Run("without a layout", func(st *testing.T) {
			r := require.New(st)

			h := e.HTML("test.html")
			r.Equal("text/html; charset=utf-8", h.ContentType())
			bb := &bytes.Buffer{}

			err := h.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("Mark", strings.TrimSpace(bb.String()))
		})

	})

	r.NoError(err)

	// t.Run("with a layout", func(st *testing.T) {
	// 	r := require.New(st)
	//
	// 	layout, err := os.Create(filepath.Join(tmpDir, "layout.html"))
	// 	r.NoError(err)
	// 	defer os.Remove(layout.Name())
	//
	// 	_, err = layout.Write([]byte("<body><%= yield %></body>"))
	// 	r.NoError(err)
	//
	// 	re := render.New(render.Options{
	// 		TemplatesBox: packr.New(tmpDir, tmpDir),
	// 		HTMLLayout:   filepath.Base(layout.Name()),
	// 	})
	//
	// 	st.Run("using just the HTMLLayout", func(sst *testing.T) {
	// 		r := require.New(sst)
	// 		h := re.HTML(filepath.Base(tmpFile.Name()))
	//
	// 		r.Equal("text/html; charset=utf-8", h.ContentType())
	// 		bb := &bytes.Buffer{}
	// 		err = h.Render(bb, map[string]interface{}{"name": "Mark"})
	// 		r.NoError(err)
	// 		r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
	// 	})
	//
	// 	st.Run("overriding the HTMLLayout", func(sst *testing.T) {
	// 		r := require.New(sst)
	// 		nlayout, err := os.Create(filepath.Join(tmpDir, "layout2.html"))
	// 		r.NoError(err)
	// 		defer os.Remove(nlayout.Name())
	//
	// 		_, err = nlayout.Write([]byte("<html><%= yield %></html>"))
	// 		r.NoError(err)
	// 		h := re.HTML(filepath.Base(tmpFile.Name()), filepath.Base(nlayout.Name()))
	//
	// 		r.Equal("text/html; charset=utf-8", h.ContentType())
	// 		bb := &bytes.Buffer{}
	// 		err = h.Render(bb, map[string]interface{}{"name": "Mark"})
	// 		r.NoError(err)
	// 		r.Equal("<html>Mark</html>", strings.TrimSpace(bb.String()))
	// 	})
	//
	// })
}
