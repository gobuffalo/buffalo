package packr

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HTTPBox(t *testing.T) {
	r := require.New(t)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(testBox))

	req, err := http.NewRequest("GET", "/hello.txt", nil)
	r.NoError(err)

	res := httptest.NewRecorder()

	mux.ServeHTTP(res, req)

	r.Equal(200, res.Code)
	r.Equal("hello world!\n", res.Body.String())
}
