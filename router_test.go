package buffalo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/markbates/going/willy"
	"github.com/stretchr/testify/require"
)

func testApp() *App {
	a := New(Options{})
	rt := a.Group("/router/tests")

	h := func(c Context) error {
		return c.Render(200, render.String(c.Request().Method))
	}

	rt.GET("/", h)
	rt.POST("/", h)
	rt.PUT("/", h)
	rt.DELETE("/", h)
	rt.OPTIONS("/", h)
	rt.PATCH("/", h)
	return a
}

func Test_Router(t *testing.T) {
	r := require.New(t)

	table := []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"OPTIONS",
		"PATCH",
	}

	ts := httptest.NewServer(testApp())
	defer ts.Close()

	for _, v := range table {
		req, err := http.NewRequest(v, fmt.Sprintf("%s/router/tests", ts.URL), nil)
		r.NoError(err)
		res, err := http.DefaultClient.Do(req)
		r.NoError(err)
		b, _ := ioutil.ReadAll(res.Body)
		r.Equal(v, string(b))
	}
}

func Test_Router_Group(t *testing.T) {
	r := require.New(t)

	a := testApp()
	g := a.Group("/api/v1")
	g.GET("/users", func(c Context) error {
		return c.NoContent(201)
	})

	w := willy.New(a)
	res := w.Request("/api/v1/users").Get()
	r.Equal(201, res.Code)
}

func Test_Router_Group_Middleware(t *testing.T) {
	r := require.New(t)

	a := testApp()
	a.Use(func(h Handler) Handler { return h })
	r.Len(a.middlewareStack.stack, 1)

	g := a.Group("/api/v1")
	r.Len(a.middlewareStack.stack, 1)
	r.Len(g.middlewareStack.stack, 1)

	g.Use(func(h Handler) Handler { return h })
	r.Len(a.middlewareStack.stack, 1)
	r.Len(g.middlewareStack.stack, 2)
}

func Test_Router_ServeFiles(t *testing.T) {
	r := require.New(t)

	tmpFile, err := ioutil.TempFile("", "assets")
	r.NoError(err)

	af := []byte("hi")
	_, err = tmpFile.Write(af)
	r.NoError(err)

	a := New(Options{})
	a.ServeFiles("/assets", http.Dir(filepath.Dir(tmpFile.Name())))

	w := willy.New(a)
	res := w.Request("/assets/%s", filepath.Base(tmpFile.Name())).Get()

	r.Equal(200, res.Code)
	r.Equal(af, res.Body.Bytes())
}
