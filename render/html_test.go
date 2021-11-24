package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/require"
)

const htmlLayout = "layout.html"
const htmlAltLayout = "alt_layout.plush.html"
const htmlTemplate = "my-template.html"

func Test_HTML_WithoutLayout(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(htmlTemplate, []byte("<%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	h := e.HTML(htmlTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("Mark", strings.TrimSpace(bb.String()))
}

func Test_HTML_WithLayout(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(htmlTemplate, []byte("<%= name %>"), 0644))
	r.NoError(rootFS.WriteFile(htmlLayout, []byte("<body><%= yield %></body>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS
	e.HTMLLayout = htmlLayout

	h := e.HTML(htmlTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
}

func Test_HTML_WithLayout_Override(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(htmlTemplate, []byte("<%= name %>"), 0644))
	r.NoError(rootFS.WriteFile(htmlLayout, []byte("<body><%= yield %></body>"), 0644))
	r.NoError(rootFS.WriteFile(htmlAltLayout, []byte("<html><%= yield %></html>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS
	e.HTMLLayout = htmlLayout

	h := e.HTML(htmlTemplate, htmlAltLayout)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<html>Mark</html>", strings.TrimSpace(bb.String()))
}
