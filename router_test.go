package buffalo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func testApp() *App {
	a := New(Options{})
	a.Redirect(301, "/foo", "/bar")
	a.GET("/bar", func(c Context) error {
		return c.Render(200, render.String("bar"))
	})

	rt := a.Group("/router/tests")

	h := func(c Context) error {
		return c.Render(200, render.String(c.Request().Method+"|"+c.Value("current_path").(string)))
	}

	rt.GET("/", h)
	rt.POST("/", h)
	rt.PUT("/", h)
	rt.DELETE("/", h)
	rt.OPTIONS("/", h)
	rt.PATCH("/", h)
	return a
}

func otherTestApp() *App {
	a := New(Options{})
	f := func(c Context) error {
		req := c.Request()
		return c.Render(200, render.String(req.Method+" - "+req.URL.String()))
	}
	a.GET("/foo", f)
	a.POST("/bar", f)
	a.DELETE("/baz/baz", f)
	return a
}

func Test_Mount_Buffalo(t *testing.T) {
	r := require.New(t)
	a := testApp()
	a.Mount("/admin", otherTestApp())

	table := map[string]string{
		"/foo":     "GET",
		"/bar":     "POST",
		"/baz/baz": "DELETE",
	}
	ts := httptest.NewServer(a)
	defer ts.Close()

	for u, m := range table {
		p := fmt.Sprintf("%s/%s", ts.URL, path.Join("admin", u))
		req, err := http.NewRequest(m, p, nil)
		r.NoError(err)
		res, err := http.DefaultClient.Do(req)
		r.NoError(err)
		b, _ := ioutil.ReadAll(res.Body)
		r.Equal(fmt.Sprintf("%s - %s", m, u), string(b))
	}
}

func Test_Mount_Buffalo_on_Group(t *testing.T) {
	r := require.New(t)
	a := testApp()
	g := a.Group("/users")
	g.Mount("/admin", otherTestApp())

	table := map[string]string{
		"/foo":     "GET",
		"/bar":     "POST",
		"/baz/baz": "DELETE",
	}
	ts := httptest.NewServer(a)
	defer ts.Close()

	for u, m := range table {
		p := fmt.Sprintf("%s/%s", ts.URL, path.Join("users", "admin", u))
		req, err := http.NewRequest(m, p, nil)
		r.NoError(err)
		res, err := http.DefaultClient.Do(req)
		r.NoError(err)
		b, _ := ioutil.ReadAll(res.Body)
		r.Equal(fmt.Sprintf("%s - %s", m, u), string(b))
	}
}

func muxer() http.Handler {
	f := func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "%s - %s", req.Method, req.URL.String())
	}
	mux := mux.NewRouter()
	mux.HandleFunc("/foo", f).Methods("GET")
	mux.HandleFunc("/bar", f).Methods("POST")
	mux.HandleFunc("/baz/baz", f).Methods("DELETE")
	return mux
}

func Test_Mount_Handler(t *testing.T) {
	r := require.New(t)
	a := testApp()
	a.Mount("/admin", muxer())

	table := map[string]string{
		"/foo":     "GET",
		"/bar":     "POST",
		"/baz/baz": "DELETE",
	}
	ts := httptest.NewServer(a)
	defer ts.Close()

	for u, m := range table {
		p := fmt.Sprintf("%s/%s", ts.URL, path.Join("admin", u))
		req, err := http.NewRequest(m, p, nil)
		r.NoError(err)
		res, err := http.DefaultClient.Do(req)
		r.NoError(err)
		b, _ := ioutil.ReadAll(res.Body)
		r.Equal(fmt.Sprintf("%s - %s", m, u), string(b))
	}
}

func Test_PreHandlers(t *testing.T) {
	r := require.New(t)
	a := testApp()
	bh := func(c Context) error {
		req := c.Request()
		return c.Render(200, render.String(req.Method+"-"+req.URL.String()))
	}
	a.GET("/ph", bh)
	a.POST("/ph", bh)
	mh := func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			res.WriteHeader(418)
			res.Write([]byte("boo"))
		}
	}
	a.PreHandlers = append(a.PreHandlers, http.HandlerFunc(mh))

	ts := httptest.NewServer(a)
	defer ts.Close()

	table := []struct {
		Code   int
		Method string
		Result string
	}{
		{Code: 418, Method: "GET", Result: "boo"},
		{Code: 200, Method: "POST", Result: "POST-/ph"},
	}

	for _, v := range table {
		req, err := http.NewRequest(v.Method, ts.URL+"/ph", nil)
		r.NoError(err)
		res, err := http.DefaultClient.Do(req)
		r.NoError(err)
		b, err := ioutil.ReadAll(res.Body)
		r.NoError(err)
		r.Equal(v.Code, res.StatusCode)
		r.Equal(v.Result, string(b))
	}
}

func Test_PreWares(t *testing.T) {
	r := require.New(t)
	a := testApp()
	bh := func(c Context) error {
		req := c.Request()
		return c.Render(200, render.String(req.Method+"-"+req.URL.String()))
	}
	a.GET("/ph", bh)
	a.POST("/ph", bh)

	mh := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				res.WriteHeader(418)
				res.Write([]byte("boo"))
			}
		})
	}

	a.PreWares = append(a.PreWares, mh)

	ts := httptest.NewServer(a)
	defer ts.Close()

	table := []struct {
		Code   int
		Method string
		Result string
	}{
		{Code: 418, Method: "GET", Result: "boo"},
		{Code: 200, Method: "POST", Result: "POST-/ph"},
	}

	for _, v := range table {
		req, err := http.NewRequest(v.Method, ts.URL+"/ph", nil)
		r.NoError(err)
		res, err := http.DefaultClient.Do(req)
		r.NoError(err)
		b, err := ioutil.ReadAll(res.Body)
		r.NoError(err)
		r.Equal(v.Code, res.StatusCode)
		r.Equal(v.Result, string(b))
	}
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
		r.Equal(fmt.Sprintf("%s|/router/tests", v), string(b))
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

func Test_Router_Group_on_Group(t *testing.T) {
	r := require.New(t)

	a := testApp()
	g := a.Group("/api/v1")
	g.GET("/users", func(c Context) error {
		return c.Render(201, nil)
	})
	f := g.Group("/foo")
	f.GET("/bar", func(c Context) error {
		return c.Render(420, nil)
	})

	w := willie.New(a)
	res := w.Request("/api/v1/foo/bar").Get()
	r.Equal(420, res.Code)
}

func Test_Router_Group_Middleware(t *testing.T) {
	r := require.New(t)

	a := testApp()
	a.Use(func(h Handler) Handler { return h })
	r.Len(a.Middleware.stack, 4)

	g := a.Group("/api/v1")
	r.Len(a.Middleware.stack, 4)
	r.Len(g.Middleware.stack, 4)

	g.Use(func(h Handler) Handler { return h })
	r.Len(a.Middleware.stack, 4)
	r.Len(g.Middleware.stack, 5)
}

func Test_Router_Redirect(t *testing.T) {
	r := require.New(t)
	w := willie.New(testApp())
	res := w.Request("/foo").Get()
	r.Equal(301, res.Code)
	r.Equal("/bar", res.Location())
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

func Test_App_NamedRoutes(t *testing.T) {

	type CarsResource struct {
		*BaseResource
	}

	r := require.New(t)
	a := New(Options{})

	var carsResource Resource
	carsResource = CarsResource{&BaseResource{}}

	rr := render.New(render.Options{
		HTMLLayout:   "application.html",
		TemplatesBox: packr.NewBox("../templates"),
		Helpers:      map[string]interface{}{},
	})

	sampleHandler := func(c Context) error {
		c.Set("opts", map[string]interface{}{})
		return c.Render(200, rr.String(`
			1. <%= rootPath() %>
			2. <%= usersPath() %>
			3. <%= userPath({user_id: 1}) %>
			4. <%= myPeepsPath() %>
			5. <%= userPath(opts) %>
			6. <%= carPath({car_id: 1}) %>
			7. <%= newCarPath() %>
			8. <%= editCarPath({car_id: 1}) %>
			9. <%= editCarPath({car_id: 1, other: 12}) %>
			10. <%= rootPath({"some":"variable","other": 12}) %>
			11. <%= rootPath() %>
			12. <%= rootPath({"special/":"12=ss"}) %>
		`))
	}

	a.GET("/", sampleHandler)
	a.GET("/users", sampleHandler)
	a.GET("/users/{user_id}", sampleHandler)
	a.GET("/peeps", sampleHandler).Name("myPeeps")
	a.Resource("/car", carsResource)

	w := willie.New(a)
	res := w.Request("/").Get()

	r.Equal(200, res.Code)
	r.Contains(res.Body.String(), "1. /")
	r.Contains(res.Body.String(), "2. /users")
	r.Contains(res.Body.String(), "3. /users/1")
	r.Contains(res.Body.String(), "4. /peeps")
	r.Contains(res.Body.String(), "5. /users/{user_id}")
	r.Contains(res.Body.String(), "6. /car/1")
	r.Contains(res.Body.String(), "7. /car/new")
	r.Contains(res.Body.String(), "8. /car/1/edit")
	r.Contains(res.Body.String(), "9. /car/1/edit?other=12")
	r.Contains(res.Body.String(), "10. /?other=12&some=variable")
	r.Contains(res.Body.String(), "11. /")
	r.Contains(res.Body.String(), "12. /?special%2F=12%3Dss")
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
			Path:   "",
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
			Path:   "",
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

	a := New(Options{})
	a.Resource("/users", &userResource{})
	a.Resource("/api/v1/users", &userResource{})

	ts := httptest.NewServer(a)
	defer ts.Close()

	c := http.Client{}
	for _, path := range []string{"/users", "/api/v1/users"} {
		for _, test := range tests {
			u := ts.URL + path + test.Path
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
	return c.Render(200, render.String(`show <%=params["user_id"] %>`))
}

func (u *userResource) New(c Context) error {
	return c.Render(200, render.String("new"))
}

func (u *userResource) Create(c Context) error {
	return c.Render(200, render.String("create"))
}

func (u *userResource) Edit(c Context) error {
	return c.Render(200, render.String(`edit <%=params["user_id"] %>`))
}

func (u *userResource) Update(c Context) error {
	return c.Render(200, render.String(`update <%=params["user_id"] %>`))
}

func (u *userResource) Destroy(c Context) error {
	return c.Render(200, render.String(`destroy <%=params["user_id"] %>`))
}

func Test_ResourceOnResource(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	ur := a.Resource("/users", &userResource{})
	ur.Resource("/people", &userResource{})

	ts := httptest.NewServer(a)
	defer ts.Close()

	type trs struct {
		Method string
		Path   string
		Result string
	}
	tests := []trs{
		{
			Method: "GET",
			Path:   "/people",
			Result: "list",
		},
		{
			Method: "GET",
			Path:   "/people/new",
			Result: "new",
		},
		{
			Method: "GET",
			Path:   "/people/1",
			Result: "show 1",
		},
		{
			Method: "GET",
			Path:   "/people/1/edit",
			Result: "edit 1",
		},
		{
			Method: "POST",
			Path:   "/people",
			Result: "create",
		},
		{
			Method: "PUT",
			Path:   "/people/1",
			Result: "update 1",
		},
		{
			Method: "DELETE",
			Path:   "/people/1",
			Result: "destroy 1",
		},
	}
	c := http.Client{}
	for _, test := range tests {
		u := ts.URL + path.Join("/users/42", test.Path)
		req, err := http.NewRequest(test.Method, u, nil)
		r.NoError(err)
		res, err := c.Do(req)
		r.NoError(err)
		b, err := ioutil.ReadAll(res.Body)
		r.NoError(err)
		r.Equal(test.Result, string(b))
	}

}

func Test_buildRouteName(t *testing.T) {
	r := require.New(t)
	cases := map[string]string{
		"/":                                          "root",
		"/users":                                     "users",
		"/users/new":                                 "newUsers",
		"/users/{user_id}":                           "user",
		"/users/{user_id}/children":                  "userChildren",
		"/users/{user_id}/children/{child_id}":       "userChild",
		"/users/{user_id}/children/new":              "newUserChildren",
		"/users/{user_id}/children/{child_id}/build": "userChildBuild",
		"/admin/planes":                              "adminPlanes",
		"/admin/planes/{plane_id}":                   "adminPlane",
		"/admin/planes/{plane_id}/edit":              "editAdminPlane",
	}

	a := New(Options{})

	for input, result := range cases {
		fResult := a.buildRouteName(input)
		r.Equal(result, fResult, input)
	}

	a = New(Options{Prefix: "/test"})
	cases = map[string]string{
		"/test":       "test",
		"/test/users": "testUsers",
	}

	for input, result := range cases {
		fResult := a.buildRouteName(input)
		r.Equal(result, fResult, input)
	}
}

func Test_CatchAll_Route(t *testing.T) {
	r := require.New(t)
	rr := render.New(render.Options{})

	a := New(Options{})
	a.GET("/{name:.+}", func(c Context) error {
		name := c.Param("name")
		return c.Render(200, rr.String(name))
	})

	w := willie.New(a)
	res := w.Request("/john").Get()

	r.Contains(res.Body.String(), "john")
}

func Test_Router_Matches_Trailing_Slash(t *testing.T) {
	r := require.New(t)

	table := []string{
		"/bar",
		"/bar/",
	}

	ts := httptest.NewServer(testApp())
	defer ts.Close()

	for _, v := range table {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", ts.URL, v), nil)
		r.NoError(err)
		res, err := http.DefaultClient.Do(req)
		r.NoError(err)
		b, _ := ioutil.ReadAll(res.Body)
		r.Equal("bar", string(b))
	}
}
