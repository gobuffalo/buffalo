package render

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"github.com/stretchr/testify/require"
)

const jsLayout = "layout.js"
const jsAltLayout = "alt_layout.plush.js"
const jsTemplate = "my-template.js"

func Test_JavaScript_WithoutLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()

	box := e.TemplatesBox
	r.NoError(box.AddString(jsTemplate, "alert(<%= name %>)"))

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("alert(Mark)", strings.TrimSpace(bb.String()))
}

// func Test_HTML_WithLayout(t *testing.T) {
// 	r := require.New(t)
//
// 	e := NewEngine()
// 	e.HTMLLayout = htmlLayout
//
// 	box := e.TemplatesBox
// 	r.NoError(box.AddString(htmlTemplate, "<%= name %>"))
// 	r.NoError(box.AddString(htmlLayout, "<body><%= yield %></body>"))
//
// 	h := e.HTML(htmlTemplate)
// 	r.Equal("text/html; charset=utf-8", h.ContentType())
// 	bb := &bytes.Buffer{}
//
// 	r.NoError(h.Render(bb, Data{"name": "Mark"}))
// 	r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
// }

func Test_JavaScript(t *testing.T) {
	r := require.New(t)

	tmpDir := filepath.Join(os.TempDir(), "markdown_test")
	err := os.MkdirAll(tmpDir, 0766)
	r.NoError(err)
	defer os.Remove(tmpDir)

	tmpFile, err := os.Create(filepath.Join(tmpDir, "test.js"))
	r.NoError(err)

	_, err = tmpFile.Write([]byte("<%= name %>"))
	r.NoError(err)

	t.Run("without a layout", func(st *testing.T) {
		r := require.New(st)

		j := New(Options{
			TemplatesBox: packr.New(tmpDir, tmpDir),
		}).JavaScript

		re := j(filepath.Base(tmpFile.Name()))
		r.Equal("application/javascript", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, map[string]interface{}{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", strings.TrimSpace(bb.String()))
	})

	t.Run("with a layout", func(st *testing.T) {
		r := require.New(st)

		layout, err := os.Create(filepath.Join(tmpDir, "layout.js"))
		r.NoError(err)

		_, err = layout.Write([]byte("<body><%= yield %></body>"))
		r.NoError(err)

		re := New(Options{
			JavaScriptLayout: filepath.Base(layout.Name()),
			TemplatesBox:     packr.New(tmpDir, tmpDir),
		})

		st.Run("using just the JavaScriptLayout", func(sst *testing.T) {
			r := require.New(sst)
			h := re.JavaScript(filepath.Base(tmpFile.Name()))

			r.Equal("application/javascript", h.ContentType())
			bb := &bytes.Buffer{}
			err = h.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
		})

		st.Run("overriding the JavaScriptLayout", func(sst *testing.T) {
			r := require.New(sst)
			nlayout, err := os.Create(filepath.Join(tmpDir, "layout2.js"))
			r.NoError(err)

			_, err = nlayout.Write([]byte("<html><%= yield %></html>"))
			r.NoError(err)
			h := re.JavaScript(filepath.Base(tmpFile.Name()), filepath.Base(nlayout.Name()))

			r.Equal("application/javascript", h.ContentType())
			bb := &bytes.Buffer{}
			err = h.Render(bb, map[string]interface{}{"name": "Mark"})
			r.NoError(err)
			r.Equal("<html>Mark</html>", strings.TrimSpace(bb.String()))
		})

	})

}

func Test_JavaScript_JS_Partial(t *testing.T) {
	r := require.New(t)

	dir, err := ioutil.TempDir("", "")
	r.NoError(err)
	defer os.RemoveAll(dir)

	re := New(Options{
		TemplatesBox: packr.New(dir, dir),
	})

	pf, err := os.Create(filepath.Join(dir, "_part.js"))
	r.NoError(err)
	_, err = pf.WriteString("alert('hi!');")
	r.NoError(err)

	tf, err := os.Create(filepath.Join(dir, "test.js"))
	r.NoError(err)
	_, err = tf.WriteString("let a = 1;\n<%= partial(\"part.js\") %>")
	r.NoError(err)

	bb := &bytes.Buffer{}
	err = re.JavaScript("test.js").Render(bb, map[string]interface{}{})
	r.NoError(err)

	r.Equal("let a = 1;\nalert('hi!');", bb.String())
}

func Test_JavaScript_JS_Partial_Without_Extension(t *testing.T) {
	r := require.New(t)

	const testJS = "let a = 1;\n<%= partial(\"part\") %>"
	const partJS = "alert('Hi <%= name %>!');"
	box := packd.NewMemoryBox()

	r.NoError(box.AddString("_part.js", partJS))
	r.NoError(box.AddString("test.js", testJS))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}

	err := re.JavaScript("test.js").Render(bb, Data{"name": "Yonghwan"})
	r.NoError(err)
	r.Equal("let a = 1;\nalert('Hi Yonghwan!');", bb.String())
}

func Test_JavaScript_HTML_Partial(t *testing.T) {
	r := require.New(t)

	const h = `<div id="foo">
	<p>hi</p>
</div>`
	box := packd.NewMemoryBox()
	r.NoError(box.AddString("_part.html", h))
	r.NoError(box.AddString("test.js", "let a = \"<%= partial(\"part.html\") %>\""))

	re := New(Options{
		TemplatesBox: box,
	})

	bb := &bytes.Buffer{}
	err := re.JavaScript("test.js").Render(bb, map[string]interface{}{})
	r.NoError(err)

	r.Equal("let a = \"\\x3Cdiv id=\\\"foo\\\"\\x3E\\u000A\\u0009\\x3Cp\\x3Ehi\\x3C/p\\x3E\\u000A\\x3C/div\\x3E\"", bb.String())
}
