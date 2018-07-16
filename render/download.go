package render

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
)

type downloadRenderer struct {
	data     []byte
	filename string
	writer   http.ResponseWriter
}

func (r downloadRenderer) ContentType() string {
	ext := filepath.Ext(r.filename)
	t := mime.TypeByExtension(ext)
	if t == "" {
		t = "application/octet-stream"
	}

	return t
}

func (r downloadRenderer) Render(w io.Writer, d Data) error {
	r.writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", r.filename))
	r.writer.Header().Add("Content-Length", strconv.Itoa(len(r.data)))

	_, err := w.Write(r.data)
	return err
}

// Download renders a file attachment automatically setting following headers:
//
//   Content-Type
//   Content-Length
//   Content-Disposition
//
// Content-Type is set using mime#TypeByExtension with the filename's extension. Content-Type will default to
// application/octet-stream if using a filename with an unknown extension.
func Download(data []byte, filename string, writer http.ResponseWriter) Renderer {
	return downloadRenderer{
		data:     data,
		filename: filename,
		writer:   writer,
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
func (e *Engine) Download(data []byte, filename string, writer http.ResponseWriter) Renderer {
	return Download(data, filename, writer)
}
