package render_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_Template(t *testing.T) {
	r := require.New(t)

	tmpFile, err := ioutil.TempFile("", "test")
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("{{name}}"))
	r.NoError(err)

	type ji func(string, ...string) render.Renderer

	table := []ji{
		render.Template,
		render.New(render.Options{}).Template,
	}

	for _, j := range table {
		re := j("foo/bar", tmpFile.Name())
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

	_, err = partFile.Write([]byte("Foo -> {{name}}"))
	r.NoError(err)

	tmpFile, err := os.Create(filepath.Join(tPath, "index.html"))
	r.NoError(err)

	_, err = tmpFile.Write([]byte(`{{partial "foo.html"}}`))
	r.NoError(err)

	type ji func(string, ...string) render.Renderer

	table := []ji{
		render.New(render.Options{TemplatesPath: tPath}).Template,
	}

	for _, j := range table {
		re := j("foo/bar", "index.html")
		r.Equal("foo/bar", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, render.Data{"name": "Mark"})
		r.NoError(err)
		r.Equal("Foo -> Mark", strings.TrimSpace(bb.String()))
	}
}

func Test_Template_WithCaching(t *testing.T) {
	r := require.New(t)

	tmpFile, err := ioutil.TempFile("", "test")
	r.NoError(err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("{{name}}"))
	r.NoError(err)

	type ji func(string, ...string) render.Renderer

	table := []ji{
		render.Template,
		render.New(render.Options{
			CacheTemplates: true,
		}).Template,
	}

	for _, j := range table {
		re := j("foo/bar", tmpFile.Name())
		r.Equal("foo/bar", re.ContentType())
		bb := &bytes.Buffer{}
		err = re.Render(bb, render.Data{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", strings.TrimSpace(bb.String()))
	}
}
