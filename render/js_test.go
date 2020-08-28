package render

import (
	"bytes"
	"strings"
	"testing"

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

func Test_JavaScript_WithLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	e.JavaScriptLayout = jsLayout

	box := e.TemplatesBox
	r.NoError(box.AddString(jsTemplate, "alert(<%= name %>)"))
	r.NoError(box.AddString(jsLayout, "$(<%= yield %>)"))

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("$(alert(Mark))", strings.TrimSpace(bb.String()))
}

func Test_JavaScript_WithLayout_Override(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	e.JavaScriptLayout = jsLayout

	box := e.TemplatesBox
	r.NoError(box.AddString(jsTemplate, "alert(<%= name %>)"))
	r.NoError(box.AddString(jsLayout, "$(<%= yield %>)"))
	r.NoError(box.AddString(jsAltLayout, "_(<%= yield %>)"))

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

	e := NewEngine()

	box := e.TemplatesBox
	r.NoError(box.AddString(jsTemplate, tmpl))
	r.NoError(box.AddString("_part.js", part))

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

	e := NewEngine()

	box := e.TemplatesBox
	r.NoError(box.AddString(jsTemplate, tmpl))
	r.NoError(box.AddString("_part.js", part))

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Yonghwan"}))
	r.Equal("let a = 1;\nalert('Hi Yonghwan!');", bb.String())
}

func Test_JavaScript_HTML_Partial(t *testing.T) {
	r := require.New(t)

	const tmpl = `let a = "<%= partial("part.html") %>"`
	const part = `<div id="foo"><p>hi</p></div>`

	e := NewEngine()
	box := e.TemplatesBox

	r.NoError(box.AddString(jsTemplate, tmpl))
	r.NoError(box.AddString("_part.html", part))

	h := e.JavaScript(jsTemplate)
	r.Equal("application/javascript", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{}))
	pre := `let a =`
	r.True(strings.HasPrefix(bb.String(), pre))
}
