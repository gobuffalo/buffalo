package render_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_DownloadWithWellKnownExtension(t *testing.T) {
	assert := require.New(t)
	data := []byte("data")

	type di func([]byte, string, http.ResponseWriter) render.Renderer
	table := []di{
		render.Download,
		render.New(render.Options{}).Download,
	}

	for _, d := range table {
		recorder := httptest.NewRecorder()
		re := d(data, "filename.pdf", recorder)
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		assert.NoError(err)
		assert.Equal(data, bb.Bytes())
		assert.Equal(strconv.Itoa(len(data)), recorder.Header().Get("Content-Length"))
		assert.Equal("attachment; filename=filename.pdf", recorder.Header().Get("Content-Disposition"))
		assert.Equal("application/pdf", re.ContentType())
	}
}

func Test_DownloadWithUnknownExtension(t *testing.T) {
	assert := require.New(t)
	data := []byte("data")

	type di func([]byte, string, http.ResponseWriter) render.Renderer
	table := []di{
		render.Download,
		render.New(render.Options{}).Download,
	}

	for _, d := range table {
		recorder := httptest.NewRecorder()
		re := d(data, "filename", recorder)
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		assert.NoError(err)
		assert.Equal(data, bb.Bytes())
		assert.Equal(strconv.Itoa(len(data)), recorder.Header().Get("Content-Length"))
		assert.Equal("attachment; filename=filename", recorder.Header().Get("Content-Disposition"))
		assert.Equal("application/octet-stream", re.ContentType())
	}
}
