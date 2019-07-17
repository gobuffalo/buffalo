package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Template(t *testing.T) {
	r := require.New(t)

	e := NewEngine()
	box := e.TemplatesBox
	r.NoError(box.AddString(htmlTemplate, `<%= name %>`))

	re := e.Template("foo/bar", htmlTemplate)
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Mark"}))
}

func Test_AssetPath(t *testing.T) {
	r := require.New(t)

	e := NewEngine()

	abox := e.AssetsBox
	r.NoError(abox.AddString("manifest.json", `{
		"application.css": "application.aabbc123.css"
	}`))

	cases := map[string]string{
		"something.txt":         "/assets/something.txt",
		"images/something.png":  "/assets/images/something.png",
		"/images/something.png": "/assets/images/something.png",
		"application.css":       "/assets/application.aabbc123.css",
	}

	for original, expected := range cases {
		tbox := e.TemplatesBox
		r.NoError(tbox.AddString(htmlTemplate, "<%= assetPath(\""+original+"\") %>"))

		re := e.Template("text/html; charset=utf-8", htmlTemplate)

		bb := &bytes.Buffer{}
		r.NoError(re.Render(bb, Data{}))
		r.Equal(expected, strings.TrimSpace(bb.String()))
	}

}

func Test_AssetPathNoManifest(t *testing.T) {
	r := require.New(t)

	e := NewEngine()

	cases := map[string]string{
		"something.txt": "/assets/something.txt",
	}

	for original, expected := range cases {
		tbox := e.TemplatesBox
		r.NoError(tbox.AddString(htmlTemplate, "<%= assetPath(\""+original+"\") %>"))

		re := e.Template("text/html; charset=utf-8", htmlTemplate)

		bb := &bytes.Buffer{}
		r.NoError(re.Render(bb, Data{}))
		r.Equal(expected, strings.TrimSpace(bb.String()))
	}
}

func Test_AssetPathNoManifestCorrupt(t *testing.T) {
	r := require.New(t)

	e := NewEngine()

	abox := e.AssetsBox
	r.NoError(abox.AddString("manifest.json", "//shdnn Corrupt!"))

	cases := map[string]string{
		"something.txt": "manifest.json is not correct",
		"other.txt":     "manifest.json is not correct",
	}

	for original, expected := range cases {
		tbox := e.TemplatesBox
		r.NoError(tbox.AddString(htmlTemplate, "<%= assetPath(\""+original+"\") %>"))

		re := e.Template("text/html; charset=utf-8", htmlTemplate)

		bb := &bytes.Buffer{}
		r.Error(re.Render(bb, Data{}))
		r.NotEqual(expected, strings.TrimSpace(bb.String()))
	}
}
