package buffalo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/markbates/willie"
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
		return c.Render(201, nil)
	})

	w := willie.New(a)
	res := w.Request("/api/v1/users").Get()
	r.Equal(201, res.Code)
}

func Test_Router_Group_Middleware(t *testing.T) {
	r := require.New(t)

	a := testApp()
	a.Use(func(h Handler) Handler { return h })
	r.Len(a.Middleware.stack, 1)

	g := a.Group("/api/v1")
	r.Len(a.Middleware.stack, 1)
	r.Len(g.Middleware.stack, 1)

	g.Use(func(h Handler) Handler { return h })
	r.Len(a.Middleware.stack, 1)
	r.Len(g.Middleware.stack, 2)
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

	w := willie.New(a)
	res := w.Request("/assets/%s", filepath.Base(tmpFile.Name())).Get()

	r.Equal(200, res.Code)
	r.Equal(af, res.Body.Bytes())
}

func Test_Resource(t *testing.T) {
	r := require.New(t)

	type trs struct {
		Method string
		Path   string
		Result string
	}

	tests := []trs{
		{
			Method: "GET",
			Path:   "/",
			Result: "list",
		},
		{
			Method: "GET",
			Path:   "/new",
			Result: "new",
		},
		{
			Method: "GET",
			Path:   "/1",
			Result: "show 1",
		},
		{
			Method: "GET",
			Path:   "/1/edit",
			Result: "edit 1",
		},
		{
			Method: "POST",
			Path:   "/",
			Result: "create",
		},
		{
			Method: "PUT",
			Path:   "/1",
			Result: "update 1",
		},
		{
			Method: "DELETE",
			Path:   "/1",
			Result: "destroy 1",
		},
	}

	a := Automatic(Options{})
	a.Resource("/users", &userResource{})
	a.Resource("/api/v1/people", &userResource{})

	ts := httptest.NewServer(a)
	defer ts.Close()

	c := http.Client{}
	for _, path := range []string{"/users", "/api/v1/people"} {
		for _, test := range tests {
			u := ts.URL + filepath.Join(path, test.Path)
			req, err := http.NewRequest(test.Method, u, nil)
			r.NoError(err)
			res, err := c.Do(req)
			r.NoError(err)
			b, err := ioutil.ReadAll(res.Body)
			r.NoError(err)
			r.Equal(test.Result, string(b))
		}
	}

}

type userResource struct{}

func (u *userResource) List(c Context) error {
	return c.Render(200, render.String("list"))
}

func (u *userResource) Show(c Context) error {
	return c.Render(200, render.String("show {{params.user_id}}{{params.person_id}}"))
}

func (u *userResource) New(c Context) error {
	return c.Render(200, render.String("new"))
}

func (u *userResource) Create(c Context) error {
	return c.Render(200, render.String("create"))
}

func (u *userResource) Edit(c Context) error {
	return c.Render(200, render.String("edit {{params.user_id}}{{params.person_id}}"))
}

func (u *userResource) Update(c Context) error {
	return c.Render(200, render.String("update {{params.user_id}}{{params.person_id}}"))
}

func (u *userResource) Destroy(c Context) error {
	return c.Render(200, render.String("destroy {{params.user_id}}{{params.person_id}}"))
}
