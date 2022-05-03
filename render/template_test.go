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

/* test if i18n files (both with plush mid-extension and latecy) proceeded correctly.
 */
func Test_Template_resolve_DefaultLang_Plush(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.plush.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"es", "en"}}))
	r.Equal("default Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_UserLang_Plush(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.plush.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko-KR", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_DefaultLang_Legacy(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"es", "en"}}))
	r.Equal("default Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_UserLang_Legacy(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko-KR", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_DefaultLang_Mixed(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	// `buffalo fix` renames templates but does not fix actions
	// in this case, aliases will be used for template matching
	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"es", "en"}}))
	r.Equal("default Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_UserLang_Mixed(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	// `buffalo fix` renames templates but does not fix actions
	// in this case, aliases will be used for template matching
	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko-KR", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

