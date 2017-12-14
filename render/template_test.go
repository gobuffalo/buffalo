package render_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/stretchr/testify/require"
)

func Test_Template(t *testing.T) {
	r := require.New(t)

	tPath, err := ioutil.TempDir("", "")
	r.NoError(err)
	defer os.Remove(tPath)

	tmpFile, err := os.Create(filepath.Join(tPath, "test"))
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("<%= name %>"))
	r.NoError(err)

	type ji func(string, ...string) render.Renderer

	table := []ji{
		render.New(render.Options{
			TemplatesBox: packr.NewBox(tPath),
		}).Template,
	}

	for _, j := range table {
		re := j("foo/bar", filepath.Base(tmpFile.Name()))
		r.Equal("foo/bar", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, render.Data{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", strings.TrimSpace(bb.String()))
	}
}

func Test_Template_Partial(t *testing.T) {
	r := require.New(t)

	tPath, err := ioutil.TempDir("", "")
	r.NoError(err)
	defer os.Remove(tPath)

	partFile, err := os.Create(filepath.Join(tPath, "_foo.html"))
	r.NoError(err)

	_, err = partFile.Write([]byte("Foo > <%= name %>"))
	r.NoError(err)

	tmpFile, err := os.Create(filepath.Join(tPath, "index.html"))
	r.NoError(err)

	_, err = tmpFile.Write([]byte(`<%= partial("foo.html") %>`))
	r.NoError(err)

	type ji func(string, ...string) render.Renderer

	j := render.New(render.Options{
		TemplatesBox: packr.NewBox(tPath),
	}).Template

	re := j("foo/bar", "index.html")
	r.Equal("foo/bar", re.ContentType())
	bb := &bytes.Buffer{}
	err = re.Render(bb, render.Data{"name": "Mark"})
	r.NoError(err)
	r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))
}

func Test_Template_Partial_WithoutExtension(t *testing.T) {
	r := require.New(t)

	tPath, err := ioutil.TempDir("", "")
	r.NoError(err)
	defer os.Remove(tPath)

	partFile, err := os.Create(filepath.Join(tPath, "_foo.html"))
	r.NoError(err)

	_, err = partFile.Write([]byte("Foo > <%= name %>"))
	r.NoError(err)

	tmpFile, err := os.Create(filepath.Join(tPath, "index.html"))
	r.NoError(err)

	_, err = tmpFile.Write([]byte(`<%= partial("foo") %>`))
	r.NoError(err)

	type ji func(string, ...string) render.Renderer

	j := render.New(render.Options{
		TemplatesBox: packr.NewBox(tPath),
	}).HTML

	re := j("index.html")
	bb := &bytes.Buffer{}
	err = re.Render(bb, render.Data{"name": "Mark"})
	r.NoError(err)
	r.Equal("Foo > Mark", strings.TrimSpace(bb.String()))
}

func Test_AssetPath(t *testing.T) {
	r := require.New(t)

	cases := map[string]string{
		"something.txt":         "/assets/something.txt",
		"images/something.png":  "/assets/images/something.png",
		"/images/something.png": "/assets/images/something.png",
		"application.css":       "/assets/application.aabbc123.css",
	}

	tDir, err := ioutil.TempDir("", "templates")
	if err != nil {
		r.Fail("Could not set the templates dir")
	}

	aDir, err := ioutil.TempDir("", "assets")
	if err != nil {
		r.Fail("Could not set the assets dir")
	}

	re := render.New(render.Options{
		TemplatesBox: packr.NewBox(tDir),
		AssetsBox:    packr.NewBox(aDir),
	}).Template

	ioutil.WriteFile(filepath.Join(aDir, "manifest.json"), []byte(`{
		"application.css": "application.aabbc123.css"
	}`), 0644)

	for original, expected := range cases {

		tmpFile, err := os.Create(filepath.Join(tDir, "test.html"))
		r.NoError(err)

		_, err = tmpFile.Write([]byte("<%= assetPath(\"" + original + "\") %>"))
		r.NoError(err)

		result := re("text/html", filepath.Base(tmpFile.Name()))

		bb := &bytes.Buffer{}
		err = result.Render(bb, render.Data{})
		r.NoError(err)
		r.Equal(expected, strings.TrimSpace(bb.String()))

		os.Remove(tmpFile.Name())
	}
}

func Test_AssetPathNoManifest(t *testing.T) {
	r := require.New(t)

	cases := map[string]string{
		"something.txt": "/assets/something.txt",
	}

	tDir, err := ioutil.TempDir("", "templates")
	if err != nil {
		r.Fail("Could not set the templates dir")
	}

	aDir, err := ioutil.TempDir("", "assets")
	if err != nil {
		r.Fail("Could not set the assets dir")
	}

	re := render.New(render.Options{
		TemplatesBox: packr.NewBox(tDir),
		AssetsBox:    packr.NewBox(aDir),
	}).Template

	for original, expected := range cases {

		tmpFile, err := os.Create(filepath.Join(tDir, "test.html"))
		r.NoError(err)

		_, err = tmpFile.Write([]byte("<%= assetPath(\"" + original + "\") %>"))
		r.NoError(err)

		result := re("text/html", filepath.Base(tmpFile.Name()))

		bb := &bytes.Buffer{}
		err = result.Render(bb, render.Data{})
		r.NoError(err)
		r.Equal(expected, strings.TrimSpace(bb.String()))

		os.Remove(tmpFile.Name())
	}
}
func Test_AssetPathManifestCorrupt(t *testing.T) {
	r := require.New(t)

	cases := map[string]string{
		"something.txt": "manifest.json is not correct",
		"other.txt":     "manifest.json is not correct",
	}

	tDir, err := ioutil.TempDir("", "templates")
	if err != nil {
		r.Fail("Could not set the templates dir")
	}

	aDir, err := ioutil.TempDir("", "assets")
	if err != nil {
		r.Fail("Could not set the assets dir")
	}

	ioutil.WriteFile(filepath.Join(aDir, "manifest.json"), []byte(`//shdnn Corrupt!`), 0644)

	re := render.New(render.Options{
		TemplatesBox: packr.NewBox(tDir),
		AssetsBox:    packr.NewBox(aDir),
	}).Template

	for original, expected := range cases {

		tmpFile, err := os.Create(filepath.Join(tDir, "test.html"))
		r.NoError(err)

		_, err = tmpFile.Write([]byte("<%= assetPath(\"" + original + "\") %>"))
		r.NoError(err)

		result := re("text/html", filepath.Base(tmpFile.Name()))

		bb := &bytes.Buffer{}
		err = result.Render(bb, render.Data{})
		r.Error(err)
		r.Contains(err.Error(), expected)

		os.Remove(tmpFile.Name())
	}
}
