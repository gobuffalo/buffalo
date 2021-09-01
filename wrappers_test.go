package buffalo

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

func Test_WrapHandlerFunc(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.GET("/foo", WrapHandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	}))

	w := httptest.New(a)
	res := w.HTML("/foo").Get()

	r.Equal("hello", res.Body.String())
}

func Test_WrapHandler(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.GET("/foo", WrapHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	})))

	w := httptest.New(a)
	res := w.HTML("/foo").Get()

	r.Equal("hello", res.Body.String())
}

func Test_WrapBuffaloHandler(t *testing.T) {
	r := require.New(t)

	tt := []struct {
		verb   string
		path   string
		status int
	}{
		{"GET", "/", 200},
		{"GET", "/foo", 201},
		{"POST", "/", 300},
		{"POST", "/foo", 400},
	}
	for _, x := range tt {
		bf := func(c Context) error {
			req := c.Request()
			return c.Render(x.status, render.String(req.Method+req.URL.Path))
		}

		h := WrapBuffaloHandler(bf)
		r.NotNil(h)

		req := httptest.NewRequest(x.verb, x.path, nil)
		res := httptest.NewRecorder()

		h.ServeHTTP(res, req)

		r.Equal(x.status, res.Code)
		r.Contains(res.Body.String(), x.verb+x.path)
	}
}

func Test_WrapBuffaloHandlerFunc(t *testing.T) {
	r := require.New(t)

	tt := []struct {
		verb   string
		path   string
		status int
	}{
		{"GET", "/", 200},
		{"GET", "/foo", 201},
		{"POST", "/", 300},
		{"POST", "/foo", 400},
	}
	for _, x := range tt {
		bf := func(c Context) error {
			req := c.Request()
			return c.Render(x.status, render.String(req.Method+req.URL.Path))
		}

		h := WrapBuffaloHandlerFunc(bf)
		r.NotNil(h)

		req := httptest.NewRequest(x.verb, x.path, nil)
		res := httptest.NewRecorder()

		h(res, req)

		r.Equal(x.status, res.Code)
		r.Contains(res.Body.String(), x.verb+x.path)
	}
}
