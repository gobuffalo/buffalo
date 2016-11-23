package buffalo

import (
	"encoding/binary"
	"net/http"
)

type buffaloResponse struct {
	status int
	size   int
	logger Logger
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
