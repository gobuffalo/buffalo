package render

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
)

type downloadRenderer struct {
	ctx    context.Context
	name   string
	reader io.Reader
}

func (r downloadRenderer) ContentType() string {
	ext := filepath.Ext(r.name)
	t := mime.TypeByExtension(ext)
	if t == "" {
		t = "application/octet-stream"
	}

	return t
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
	header.Add("Content-Disposition", disposition)
	contentLength := strconv.Itoa(int(written))
	header.Add("Content-Length", contentLength)

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
//
// Note: the purpose of this function is not serving static files but to support
// downloading of dynamically genrated data as a file. For example, you can use
// this function when you implement CSV file download feature for the result of
// a database query.
//
// Do not use this function for large io.Reader. It could cause out of memory if
// the size of io.Reader is too big.
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
//
// Note: the purpose of this method is not serving static files but to support
// downloading of dynamically genrated data as a file. For example, you can use
// this method when you implement CSV file download feature for the result of
// a database query.
//
// Do not use this method for large io.Reader. It could cause out of memory if
// the size of io.Reader is too big.
func (e *Engine) Download(ctx context.Context, name string, r io.Reader) Renderer {
	return Download(ctx, name, r)
}

type responsible interface {
	Response() http.ResponseWriter
}
