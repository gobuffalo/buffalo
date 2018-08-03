package render_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/stretchr/testify/require"
)

var data = []byte("data")

func Test_Download_KnownExtension(t *testing.T) {
	assert := require.New(t)

	type di func(context.Context, string, io.Reader) render.Renderer
	table := []di{
		render.Download,
		render.New(render.Options{}).Download,
	}

	for _, d := range table {
		ctx := testContext{rw: httptest.NewRecorder()}
		re := d(ctx, "filename.pdf", bytes.NewReader(data))
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		assert.NoError(err)
		assert.Equal(data, bb.Bytes())
		assert.Equal(strconv.Itoa(len(data)), ctx.Response().Header().Get("Content-Length"))
		assert.Equal("attachment; filename=filename.pdf", ctx.Response().Header().Get("Content-Disposition"))
		assert.Equal("application/pdf", re.ContentType())
	}
}

func Test_Download_UnknownExtension(t *testing.T) {
	assert := require.New(t)

	type di func(context.Context, string, io.Reader) render.Renderer
	table := []di{
		render.Download,
		render.New(render.Options{}).Download,
	}

	for _, d := range table {
		ctx := testContext{rw: httptest.NewRecorder()}
		re := d(ctx, "filename", bytes.NewReader(data))
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		assert.NoError(err)
		assert.Equal(data, bb.Bytes())
		assert.Equal(strconv.Itoa(len(data)), ctx.Response().Header().Get("Content-Length"))
		assert.Equal("attachment; filename=filename", ctx.Response().Header().Get("Content-Disposition"))
		assert.Equal("application/octet-stream", re.ContentType())
	}
}

func Test_InvalidContext(t *testing.T) {
	assert := require.New(t)

	type di func(context.Context, string, io.Reader) render.Renderer
	table := []di{
		render.Download,
		render.New(render.Options{}).Download,
	}

	for _, d := range table {
		ctx := context.TODO()
		re := d(ctx, "filename", bytes.NewReader(data))
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		assert.Error(err)
	}
}

type testContext struct {
	context.Context
	rw http.ResponseWriter
}

func (c testContext) Response() http.ResponseWriter {
	return c.rw
}
