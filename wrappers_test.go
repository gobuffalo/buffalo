package buffalo

import (
	"net/http"
	"testing"

	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_WrapHandlerFunc(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.GET("/foo", WrapHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	}))

	w := willie.New(a)
	res := w.Request("/foo").Get()

	r.Equal("hello", res.Body.String())
}

func Test_WrapHandler(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.GET("/foo", WrapHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	})))

	w := willie.New(a)
	res := w.Request("/foo").Get()

	r.Equal("hello", res.Body.String())
}
