package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/require"
)

const jsLayout = "layout.js"
const jsAltLayout = "alt_layout.plush.js"
const jsTemplate = "my-template.js"

func Test_JavaScript_WithoutLayout(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(jsTemplate, []byte("alert(<%= name %>)"), 0644))

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

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(jsTemplate, []byte("alert(<%= name %>)"), 0644))
	r.NoError(rootFS.WriteFile(jsLayout, []byte("$(<%= yield %>)"), 0644))

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

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(jsTemplate, []byte("alert(<%= name %>)"), 0644))
	r.NoError(rootFS.WriteFile(jsLayout, []byte("$(<%= yield %>)"), 0644))
	r.NoError(rootFS.WriteFile(jsAltLayout, []byte("_(<%= yield %>)"), 0644))

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

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(jsTemplate, []byte(tmpl), 0644))
	r.NoError(rootFS.WriteFile("_part.js", []byte(part), 0644))

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

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(jsTemplate, []byte(tmpl), 0644))
	r.NoError(rootFS.WriteFile("_part.js", []byte(part), 0644))

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

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(jsTemplate, []byte(tmpl), 0644))
	r.NoError(rootFS.WriteFile("_part.html", []byte(part), 0644))

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
