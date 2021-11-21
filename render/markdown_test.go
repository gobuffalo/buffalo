package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/require"
)

const mdTemplate = "my-template.md"

func Test_MD_WithoutLayout(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(mdTemplate, []byte("<%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	h := e.HTML(mdTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<p>Mark</p>", strings.TrimSpace(bb.String()))
}

func Test_MD_WithLayout(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(mdTemplate, []byte("<%= name %>"), 0644))
	r.NoError(rootFS.WriteFile(htmlLayout, []byte("<body><%= yield %></body>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS
	e.HTMLLayout = htmlLayout

	h := e.HTML(mdTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<body><p>Mark</p>\n</body>", strings.TrimSpace(bb.String()))
}
