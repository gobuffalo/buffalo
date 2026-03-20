package render

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

const jsLayout = "layout.js"
const jsAltLayout = "alt_layout.plush.js"
const jsTemplate = "my-template.js"

func Test_JavaScript_WithoutLayout(t *testing.T) {
	r := require.New(t)

	rootFS := fstest.MapFS{
		jsTemplate: &fstest.MapFile{
			Data: []byte("alert(<%= name %>)"),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("alert(Mark)", strings.TrimSpace(bb.String()))
}

func Test_JavaScript_WithLayout(t *testing.T) {
	r := require.New(t)

	rootFS := fstest.MapFS{
		jsTemplate: &fstest.MapFile{
			Data: []byte("alert(<%= name %>)"),
			Mode: 0644,
		},
		jsLayout: &fstest.MapFile{
			Data: []byte("$(<%= yield %>)"),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS
	e.JavaScriptLayout = jsLayout

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("$(alert(Mark))", strings.TrimSpace(bb.String()))
}

func Test_JavaScript_WithLayout_Override(t *testing.T) {
	r := require.New(t)

	rootFS := fstest.MapFS{
		jsTemplate: &fstest.MapFile{
			Data: []byte("alert(<%= name %>)"),
			Mode: 0644,
		},
		jsLayout: &fstest.MapFile{
			Data: []byte("$(<%= yield %>)"),
			Mode: 0644,
		},
		jsAltLayout: &fstest.MapFile{
			Data: []byte("_(<%= yield %>)"),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS
	e.JavaScriptLayout = jsLayout

	h := e.JavaScript(jsTemplate, jsAltLayout)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("_(alert(Mark))", strings.TrimSpace(bb.String()))
}

func Test_JavaScript_Partial_Without_Extension(t *testing.T) {
	const tmpl = "let a = 1;\n<%= partial(\"part\") %>"
	const part = "alert('Hi <%= name %>!');"

	r := require.New(t)

	rootFS := fstest.MapFS{
		jsTemplate: &fstest.MapFile{
			Data: []byte(tmpl),
			Mode: 0644,
		},
		"_part.js": &fstest.MapFile{
			Data: []byte(part),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS
	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Yonghwan"}))
	r.Equal("let a = 1;\nalert('Hi Yonghwan!');", bb.String())
}

func Test_JavaScript_Partial(t *testing.T) {
	const tmpl = "let a = 1;\n<%= partial(\"part.js\") %>"
	const part = "alert('Hi <%= name %>!');"

	r := require.New(t)

	rootFS := fstest.MapFS{
		jsTemplate: &fstest.MapFile{
			Data: []byte(tmpl),
			Mode: 0644,
		},
		"_part.js": &fstest.MapFile{
			Data: []byte(part),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Yonghwan"}))
	r.Equal("let a = 1;\nalert('Hi Yonghwan!');", bb.String())
}

func Test_JavaScript_HTML_Partial(t *testing.T) {
	const tmpl = "let a = \"<%= partial(\"part.html\") %>\""
	const part = `<div id="foo">
	<p>hi</p>
</div>`

	r := require.New(t)

	rootFS := fstest.MapFS{
		jsTemplate: &fstest.MapFile{
			Data: []byte(tmpl),
			Mode: 0644,
		},
		"_part.html": &fstest.MapFile{
			Data: []byte(part),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{}))
	r.Contains(bb.String(), `id`)
	r.Contains(bb.String(), `foo`)

	// To check it has escaped the partial
	r.NotContains(bb.String(), `<div id="foo">`)
}
