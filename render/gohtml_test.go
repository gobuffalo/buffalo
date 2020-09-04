package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const gohtmlLayout = "go-layout.tmpl"
const gohtmlTemplate = "go-template.tmpl"
const gohtmlInclude = "go-include.tmpl"

func Test_GoHTML_WithoutLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	box := e.TemplatesBox

	r.NoError(box.AddString(gohtmlTemplate, "{{ .name }}"))

	h := e.HTML(gohtmlTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("Mark", strings.TrimSpace(bb.String()))
}

func Test_GoHTML_Include(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	box := e.TemplatesBox

	r.NoError(box.AddString(gohtmlTemplate, `{{ .name }} {{ include "go-include.tmpl" }} `))
	r.NoError(box.AddString(gohtmlInclude, "blebleble"))

	h := e.HTML(gohtmlTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("Mark blebleble", strings.TrimSpace(bb.String()))
}

func Test_GoHTML_WithLayout(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	e.HTMLLayout = gohtmlLayout

	box := e.TemplatesBox
	r.NoError(box.AddString(gohtmlTemplate, "{{ .name }}"))
	r.NoError(box.AddString(gohtmlLayout, "<body>{{ .yield }}</body>"))

	h := e.HTML(gohtmlTemplate)
	r.Equal("text/html; charset=utf-8", h.ContentType())
	bb := &bytes.Buffer{}

	r.NoError(h.Render(bb, Data{"name": "Mark"}))
	r.Equal("<body>Mark</body>", strings.TrimSpace(bb.String()))
}
