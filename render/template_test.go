package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/require"
)

func Test_Template(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile(htmlTemplate, []byte("<%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", htmlTemplate)
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Mark"}))
}

func Test_AssetPath(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("manifest.json", []byte(`{
		"application.css": "application.aabbc123.css"
	}`), 0644))

	e := NewEngine()
	e.AssetsFS = rootFS

	cases := map[string]string{
		"something.txt":         "/assets/something.txt",
		"images/something.png":  "/assets/images/something.png",
		"/images/something.png": "/assets/images/something.png",
		"application.css":       "/assets/application.aabbc123.css",
	}

	for original, expected := range cases {
		rootFS := memfs.New()
		r.NoError(rootFS.WriteFile(htmlTemplate, []byte("<%= assetPath(\""+original+"\") %>"), 0644))
		e.TemplatesFS = rootFS

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
		rootFS := memfs.New()
		r.NoError(rootFS.WriteFile(htmlTemplate, []byte("<%= assetPath(\""+original+"\") %>"), 0644))
		e.TemplatesFS = rootFS

		re := e.Template("text/html; charset=utf-8", htmlTemplate)

		bb := &bytes.Buffer{}
		r.NoError(re.Render(bb, Data{}))
		r.Equal(expected, strings.TrimSpace(bb.String()))
	}
}

func Test_AssetPathNoManifestCorrupt(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("manifest.json", []byte("//shdnn Corrupt!"), 0644))

	e := NewEngine()
	e.AssetsFS = rootFS

	cases := map[string]string{
		"something.txt": "manifest.json is not correct",
		"other.txt":     "manifest.json is not correct",
	}

	for original, expected := range cases {
		rootFS := memfs.New()
		r.NoError(rootFS.WriteFile(htmlTemplate, []byte("<%= assetPath(\""+original+"\") %>"), 0644))
		e.TemplatesFS = rootFS

		re := e.Template("text/html; charset=utf-8", htmlTemplate)

		bb := &bytes.Buffer{}
		r.Error(re.Render(bb, Data{}))
		r.NotEqual(expected, strings.TrimSpace(bb.String()))
	}
}
