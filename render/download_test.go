package render

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type dlRenderer func(context.Context, string, io.Reader) Renderer

type dlContext struct {
	context.Context
	rw http.ResponseWriter
}

func (c dlContext) Response() http.ResponseWriter {
	return c.rw
}

var data = []byte("data")

func Test_Download_KnownExtension(t *testing.T) {
	r := require.New(t)

	table := []dlRenderer{
		Download,
		New(Options{}).Download,
	}

	for _, dl := range table {
		ctx := dlContext{rw: httptest.NewRecorder()}

		re := dl(ctx, "filename.pdf", bytes.NewReader(data))

		bb := &bytes.Buffer{}
		r.NoError(re.Render(bb, nil))

		r.Equal(data, bb.Bytes())
		r.Equal(strconv.Itoa(len(data)), ctx.Response().Header().Get("Content-Length"))
		r.Equal("attachment; filename=filename.pdf", ctx.Response().Header().Get("Content-Disposition"))
		r.Equal("application/pdf", re.ContentType())
	}
}

func Test_Download_UnknownExtension(t *testing.T) {
	r := require.New(t)

	table := []dlRenderer{
		Download,
		New(Options{}).Download,
	}

	for _, dl := range table {
		ctx := dlContext{rw: httptest.NewRecorder()}
		re := dl(ctx, "filename", bytes.NewReader(data))

		bb := &bytes.Buffer{}
		r.NoError(re.Render(bb, nil))

		r.Equal(data, bb.Bytes())
		r.Equal(strconv.Itoa(len(data)), ctx.Response().Header().Get("Content-Length"))
		r.Equal("attachment; filename=filename", ctx.Response().Header().Get("Content-Disposition"))
		r.Equal("application/octet-stream", re.ContentType())
	}
}

func Test_InvalidContext(t *testing.T) {
	r := require.New(t)

	table := []dlRenderer{
		Download,
		New(Options{}).Download,
	}

	for _, dl := range table {
		re := dl(context.TODO(), "filename", bytes.NewReader(data))

		bb := &bytes.Buffer{}
		r.Error(re.Render(bb, nil))
	}
}
