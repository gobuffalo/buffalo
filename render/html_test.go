package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const htmlLayout = "layout.html"
const htmlAltLayout = "alt_layout.plush.html"
const htmlTemplate = "my-template.html"

func Test_HTML_WithoutLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()

	box := e.TemplatesBox
	r.NoError(box.AddString(htmlTemplate, "<%= name %>"))

	h := e.HTML(htmlTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	data := map[string]interface{}{"name": "Mark"}
	r.NoError(h.Render(bb, data))
	r.Equal("Mark", strings.TrimSpace(bb.String()))
}

func Test_HTML_WithLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	e.HTMLLayout = htmlLayout

	box := e.TemplatesBox
	r.NoError(box.AddString(htmlTemplate, "<%= name %>"))
	r.NoError(box.AddString(htmlLayout, "<body><%= yield %></body>"))

	h := e.HTML(htmlTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	data := map[string]interface{}{"name": "Mark"}
	r.NoError(h.Render(bb, data))
	r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
}

func Test_HTML_WithLayout_Override(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	e.HTMLLayout = htmlLayout

	box := e.TemplatesBox
	r.NoError(box.AddString(htmlTemplate, "<%= name %>"))
	r.NoError(box.AddString(htmlLayout, "<body><%= yield %></body>"))
	r.NoError(box.AddString(htmlAltLayout, "<html><%= yield %></html>"))

	h := e.HTML(htmlTemplate, htmlAltLayout)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	data := map[string]interface{}{"name": "Mark"}
	r.NoError(h.Render(bb, data))
	r.Equal("<html>Mark</html>", strings.TrimSpace(bb.String()))
}
