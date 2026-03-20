package render

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

const mdTemplate = "my-template.md"

func Test_MD_WithoutLayout(t *testing.T) {
	r := require.New(t)

	rootFS := fstest.MapFS{
		mdTemplate: &fstest.MapFile{
			Data: []byte("<%= name %>"),
			Mode: 0644,
		},
	}

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

	rootFS := fstest.MapFS{
		mdTemplate: &fstest.MapFile{
			Data: []byte("<%= name %>"),
			Mode: 0644,
		},
		htmlLayout: &fstest.MapFile{
			Data: []byte("<body><%= yield %></body>"),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS
	e.HTMLLayout = htmlLayout

	h := e.HTML(mdTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<body><p>Mark</p>\n</body>", strings.TrimSpace(bb.String()))
}
