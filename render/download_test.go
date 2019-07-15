package render

import (
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

type di func(context.Context, string, io.Reader) Renderer

var data = []byte("data")

func Test_Download_KnownExtension(t *testing.T) {
	r := require.New(t)

	table := []di{
		Download,
		New(Options{}).Download,
	}

	for _, d := range table {
		ctx := testContext{rw: httptest.NewRecorder()}
		re := d(ctx, "filename.pdf", bytes.NewReader(data))
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		r.NoError(err)
		r.Equal(data, bb.Bytes())
		r.Equal(strconv.Itoa(len(data)), ctx.Response().Header().Get("Content-Length"))
		r.Equal("attachment; filename=filename.pdf", ctx.Response().Header().Get("Content-Disposition"))
		r.Equal("application/pdf", re.ContentType())
	}
}

func Test_Download_UnknownExtension(t *testing.T) {
	r := require.New(t)

	table := []di{
		Download,
		New(Options{}).Download,
	}

	for _, d := range table {
		ctx := testContext{rw: httptest.NewRecorder()}
		re := d(ctx, "filename", bytes.NewReader(data))
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		r.NoError(err)
		r.Equal(data, bb.Bytes())
		r.Equal(strconv.Itoa(len(data)), ctx.Response().Header().Get("Content-Length"))
		r.Equal("attachment; filename=filename", ctx.Response().Header().Get("Content-Disposition"))
		r.Equal("application/octet-stream", re.ContentType())
	}
}

func Test_InvalidContext(t *testing.T) {
	r := require.New(t)

	table := []di{
		Download,
		New(Options{}).Download,
	}

	for _, d := range table {
		ctx := context.TODO()
		re := d(ctx, "filename", bytes.NewReader(data))
		bb := new(bytes.Buffer)
		err := re.Render(bb, nil)

		r.Error(err)
	}
}
