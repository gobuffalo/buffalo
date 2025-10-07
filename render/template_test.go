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

// support short language-only version of template e.g. index.plush.ko.html
func Test_Template_resolve_FullLocale_ShortFile(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.plush.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko-KR", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_LangOnly_FullFile(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.plush.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_FullLocale_ShortFile_Legacy(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.ko.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko-KR", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_LangOnly_FullFile_Legacy(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_FullLocale_ShortFile_Mixed(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko-KR", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_resolve_LangOnly_FullFile_Mixed(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("index.plush.html", []byte("default <%= name %>"), 0644))
	r.NoError(rootFS.WriteFile("index.plush.ko-kr.html", []byte("korean <%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Template("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())

	bb := &bytes.Buffer{}
	r.NoError(re.Render(bb, Data{"name": "Paul", "languages": []string{"ko", "en"}}))
	r.Equal("korean Paul", strings.TrimSpace(bb.String()))
}

func Test_Template_extsAndBase(t *testing.T) {
	r := require.New(t)

	tests := []struct {
		name         string
		input        string
		expectedExts []string
		expectedBase string
	}{
		{
			name:         "single extension",
			input:        "index.html",
			expectedExts: []string{"html"},
			expectedBase: "index",
		},
		{
			name:         "multiple extensions",
			input:        "template.html.plush",
			expectedExts: []string{"plush", "html"},
			expectedBase: "template",
		},
		{
			name:         "three extensions",
			input:        "layout.md.html.plush",
			expectedExts: []string{"plush", "md", "html"},
			expectedBase: "layout",
		},
		{
			name:         "no extension",
			input:        "template",
			expectedExts: []string{"html"},
			expectedBase: "template",
		},
		{
			name:         "empty string",
			input:        "",
			expectedExts: []string{"html"},
			expectedBase: "",
		},
		{
			name:         "extension with uppercase",
			input:        "index.HTML",
			expectedExts: []string{"html"},
			expectedBase: "index",
		},
		{
			name:         "mixed case extensions",
			input:        "template.MD.HTML.PLUSH",
			expectedExts: []string{"plush", "md", "html"},
			expectedBase: "template",
		},
		{
			name:         "path with directories",
			input:        "layouts/application.html.plush",
			expectedExts: []string{"plush", "html"},
			expectedBase: "layouts/application",
		},
		{
			name:         "nested path no extension",
			input:        "views/users/index",
			expectedExts: []string{"html"},
			expectedBase: "views/users/index",
		},
		{
			name:         "dotfile",
			input:        ".gitignore",
			expectedExts: []string{"gitignore"},
			expectedBase: "",
		},
		{
			name:         "dotfile with extension",
			input:        ".env.local",
			expectedExts: []string{"local", "env"},
			expectedBase: "",
		},
		{
			name:         "complex filename",
			input:        "user-profile.en-US.html.plush",
			expectedExts: []string{"plush", "html", "en-us"},
			expectedBase: "user-profile",
		},
		{
			name:         "only extension",
			input:        ".html",
			expectedExts: []string{"html"},
			expectedBase: "",
		},
		{
			name:         "locale and template extensions",
			input:        "welcome.fr-FR.html.plush",
			expectedExts: []string{"plush", "html", "fr-fr"},
			expectedBase: "welcome",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := &templateRenderer{}

			gotExts, gotBase := renderer.extsAndBase(tt.input)

			r.Equal(tt.expectedExts, gotExts, "extensions should match")
			r.Equal(tt.expectedBase, gotBase, "base name should match")
		})
	}
}
