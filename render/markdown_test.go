package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const mdTemplate = "my-template.md"

func Test_MD_WithoutLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()

	box := e.TemplatesBox
	r.NoError(box.AddString(mdTemplate, "<%= name %>"))

	h := e.HTML(mdTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<p>Mark</p>", strings.TrimSpace(bb.String()))
}

func Test_MD_WithLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	e.HTMLLayout = htmlLayout

	box := e.TemplatesBox
	r.NoError(box.AddString(mdTemplate, "<%= name %>"))
	r.NoError(box.AddString(htmlLayout, "<body><%= yield %></body>"))

	h := e.HTML(mdTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<body><p>Mark</p>\n</body>", strings.TrimSpace(bb.String()))
}
