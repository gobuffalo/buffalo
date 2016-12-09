package buffalo

import (
	"bufio"
	"encoding/binary"
	"net"
	"net/http"

	"github.com/pkg/errors"
)

type buffaloResponse struct {
	status int
	size   int
	http.ResponseWriter
}

func (w *buffaloResponse) WriteHeader(i int) {
	w.status = i
	w.ResponseWriter.WriteHeader(i)
}

func (w *buffaloResponse) Write(b []byte) (int, error) {
	w.size = binary.Size(b)
	return w.ResponseWriter.Write(b)
}
func (w *buffaloResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.WithStack(errors.New("does not implement http.Hijack"))
}

func (w *buffaloResponse) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *buffaloResponse) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
