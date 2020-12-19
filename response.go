package buffalo

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
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
func (w *Response) WriteHeader(code int) {
	if code == w.Status {
		return
	}

	if w.Status > 0 {
		fmt.Printf("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.Status, code)
		return
	}

	w.Status = code
	w.ResponseWriter.WriteHeader(code)
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
	return nil, nil, fmt.Errorf("does not implement http.Hijack")
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
