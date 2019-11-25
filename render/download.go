package render

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gobuffalo/buffalo/internal/consts"
)

type downloadRenderer struct {
	ctx    context.Context
	name   string
	reader io.Reader
}

func (r downloadRenderer) ContentType() string {
	ext := filepath.Ext(r.name)
	t := mime.TypeByExtension(ext)

	if len(t) != 0 {
		return t
	}

	return consts.MIME_Octet_Stream
}

func (r downloadRenderer) Render(w io.Writer, d Data) error {
	written, err := io.Copy(w, r.reader)
	if err != nil {
		return err
	}

	ctx, ok := r.ctx.(responsible)
	if !ok {
		return fmt.Errorf("context has no response writer")
	}

	header := ctx.Response().Header()
	disposition := fmt.Sprintf("attachment; filename=%s", r.name)
	header.Add(consts.HTTP_ContentDisposition, disposition)
	contentLength := strconv.Itoa(int(written))
	header.Add(consts.HTTP_ContentLength, contentLength)

	return nil
}

// Download renders a file attachment automatically setting following headers:
//
//   Content-Type
//   Content-Length
//   Content-Disposition
//
// Content-Type is set using mime#TypeByExtension with the filename's extension. Content-Type will default to
// application/octet-stream if using a filename with an unknown extension.
func Download(ctx context.Context, name string, r io.Reader) Renderer {
	return downloadRenderer{
		ctx:    ctx,
		name:   name,
		reader: r,
	}
}

// Download renders a file attachment automatically setting following headers:
//
//   Content-Type
//   Content-Length
//   Content-Disposition
//
// Content-Type is set using mime#TypeByExtension with the filename's extension. Content-Type will default to
// application/octet-stream if using a filename with an unknown extension.
func (e *Engine) Download(ctx context.Context, name string, r io.Reader) Renderer {
	return Download(ctx, name, r)
}

type responsible interface {
	Response() http.ResponseWriter
}
