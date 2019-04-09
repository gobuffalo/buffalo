package buffalo

import (
	"bufio"
	"encoding/binary"
	"net"
	"net/http"

	"errors"
)

// Response implements the http.ResponseWriter interface and allows
// for the capture of the response status and size to be used for things
// like logging requests.
type Response struct {
	Status int
	Size   int
	http.ResponseWriter
}

// WriteHeader sets the status code for a response
func (w *Response) WriteHeader(i int) {
	w.Status = i
	w.ResponseWriter.WriteHeader(i)
}

// Write the body of the response
func (w *Response) Write(b []byte) (int, error) {
	w.Size = binary.Size(b)
	return w.ResponseWriter.Write(b)
}

// Hijack implements the http.Hijacker interface to allow for things like websockets.
func (w *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("does not implement http.Hijack")
}

// Flush the response
func (w *Response) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

type closeNotifier interface {
	CloseNotify() <-chan bool
}

// CloseNotify implements the http.CloseNotifier interface
func (w *Response) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(closeNotifier); ok {
		return cn.CloseNotify()
	}
	return nil
}
