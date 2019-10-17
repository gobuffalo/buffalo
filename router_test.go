package buffalo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func testApp() *App {
	a := New(Options{})
	a.Redirect(http.StatusMovedPermanently, "/foo", "/bar")
	a.GET("/bar", func(c Context) error {
		return c.Render(http.StatusOK, render.String("bar"))
	})

	rt := a.Group("/router/tests")

	h := func(c Context) error {
		x := c.Request().Method + "|"
		x += strings.TrimSuffix(c.Value("current_path").(string), "/")
		return c.Render(http.StatusOK, render.String(x))
	}

	rt.GET("/", h)
	rt.POST("/", h)
	rt.PUT("/", h)
	rt.DELETE("/", h)
	rt.OPTIONS("/", h)
	rt.PATCH("/", h)

	a.ErrorHandlers[http.StatusMethodNotAllowed] = func(status int, err error, c Context) error {
		res := c.Response()
		res.WriteHeader(status)
		res.Write([]byte("my custom 405"))
		return nil
	}
	return a
}

func otherTestApp() *App {
	a := New(Options{})
	f := func(c Context) error {
		req := c.Request()
		return c.Render(http.StatusOK, render.String(req.Method+" - "+req.URL.String()))
	}
	a.GET("/foo", f)
	a.POST("/bar", f)
	a.DELETE("/baz/baz", f)
	return a
}

func Test_MethodNotFoundError(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.GET("/bar", func(c Context) error {
		return c.Render(http.StatusOK, render.String("bar"))
	})
	a.ErrorHandlers[http.StatusMethodNotAllowed] = func(status int, err error, c Context) error {
		res := c.Response()
		res.WriteHeader(status)
		res.Write([]byte("my custom 405"))
		return nil
	}
	w := httptest.New(a)
	res := w.HTML("/bar").Post(nil)
	r.Equal(http.StatusMethodNotAllowed, res.Code)
	r.Contains(res.Body.String(), "my custom 405")
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
		r.Equal(fmt.Sprintf("%s - %s/", m, u), string(b))
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
		r.Equal(fmt.Sprintf("%s - %s/", m, u), string(b))
	}
}

func muxer() http.Handler {
	f := func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "%s - %s", req.Method, req.URL.String())
	}
	mux := mux.NewRouter()
	mux.HandleFunc("/foo/", f).Methods("GET")
	mux.HandleFunc("/bar/", f).Methods("POST")
	mux.HandleFunc("/baz/baz/", f).Methods("DELETE")
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
		r.Equal(fmt.Sprintf("%s - %s/", m, u), string(b))
	}
}

func Test_PreHandlers(t *testing.T) {
	r := require.New(t)
	a := testApp()
	bh := func(c Context) error {
		req := c.Request()
		return c.Render(http.StatusOK, render.String(req.Method+"-"+req.URL.String()))
	}
	a.GET("/ph", bh)
	a.POST("/ph", bh)
	mh := func(res http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			res.WriteHeader(http.StatusTeapot)
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
		{Code: http.StatusTeapot, Method: "GET", Result: "boo"},
		{Code: http.StatusOK, Method: "POST", Result: "POST-/ph/"},
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
		return c.Render(http.StatusOK, render.String(req.Method+"-"+req.URL.String()))
	}
	a.GET("/ph", bh)
	a.POST("/ph", bh)

	mh := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.Method == "GET" {
				res.WriteHeader(http.StatusTeapot)
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
		{Code: http.StatusTeapot, Method: "GET", Result: "boo"},
		{Code: http.StatusOK, Method: "POST", Result: "POST-/ph/"},
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
		return c.Render(http.StatusCreated, nil)
	})

	w := httptest.New(a)
	res := w.HTML("/api/v1/users").Get()
	r.Equal(http.StatusCreated, res.Code)
}

func Test_Router_Group_on_Group(t *testing.T) {
	r := require.New(t)

	a := testApp()
	g := a.Group("/api/v1")
	g.GET("/users", func(c Context) error {
		return c.Render(http.StatusCreated, nil)
	})
	f := g.Group("/foo")
	f.GET("/bar", func(c Context) error {
		return c.Render(http.StatusTeapot, nil)
	})

	w := httptest.New(a)
	res := w.HTML("/api/v1/foo/bar").Get()
	r.Equal(http.StatusTeapot, res.Code)
}

func Test_Router_Group_Middleware(t *testing.T) {
	r := require.New(t)

	a := testApp()
	a.Use(func(h Handler) Handler { return h })
	r.Len(a.Middleware.stack, 5)

	g := a.Group("/api/v1")
	r.Len(a.Middleware.stack, 5)
	r.Len(g.Middleware.stack, 5)

	g.Use(func(h Handler) Handler { return h })
	r.Len(a.Middleware.stack, 5)
	r.Len(g.Middleware.stack, 6)
}

func Test_Router_Redirect(t *testing.T) {
	r := require.New(t)
	w := httptest.New(testApp())
	res := w.HTML("/foo").Get()
	r.Equal(http.StatusMovedPermanently, res.Code)
	r.Equal("/bar", res.Location())
}

func Test_Router_ServeFiles(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	box.AddString("foo.png", "foo")
	a := New(Options{})
	a.ServeFiles("/assets", box)

	w := httptest.New(a)
	res := w.HTML("/assets/foo.png").Get()

	r.Equal(http.StatusOK, res.Code)
	r.Equal("foo", res.Body.String())

	r.NotEqual(res.Header().Get("ETag"), "")
	r.Equal(res.Header().Get("Cache-Control"), "max-age=31536000")

	envy.Set(AssetsAgeVarName, "3600")
	w = httptest.New(a)
	res = w.HTML("/assets/foo.png").Get()

	r.Equal(http.StatusOK, res.Code)
	r.Equal("foo", res.Body.String())

	r.NotEqual(res.Header().Get("ETag"), "")
	r.Equal(res.Header().Get("Cache-Control"), "max-age=3600")
}

func Test_Router_InvalidURL(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	box.AddString("foo.png", "foo")
	a := New(Options{})
	a.ServeFiles("/", box)

	w := httptest.New(a)
	s := "/%25%7dn2zq0%3cscript%3ealert(1)%3c\\/script%3evea7f"

	request, _ := http.NewRequest("GET", s, nil)
	response := httptest.NewRecorder()

	w.ServeHTTP(response, request)
	r.Equal(http.StatusBadRequest, response.Code, "(400) BadRequest response is expected")
}

type WebResource struct {
	BaseResource
}

// Edit default implementation. Returns a 404
func (v WebResource) Edit(c Context) error {
	return c.Error(http.StatusNotFound, fmt.Errorf("resource not implemented"))
}

// New default implementation. Returns a 404
func (v WebResource) New(c Context) error {
	return c.Error(http.StatusNotFound, fmt.Errorf("resource not implemented"))
}

func Test_App_NamedRoutes(t *testing.T) {

	type CarsResource struct {
		WebResource
	}

	type ResourcesResource struct {
		WebResource
	}

	r := require.New(t)
	a := New(Options{})

	var carsResource Resource = CarsResource{}

	var resourcesResource Resource = ResourcesResource{}

	rr := render.New(render.Options{
		HTMLLayout:   "application.plush.html",
		TemplatesBox: packr.New("../templates", "../templates"),
		Helpers:      map[string]interface{}{},
	})

	sampleHandler := func(c Context) error {
		c.Set("opts", map[string]interface{}{})
		return c.Render(http.StatusOK, rr.String(`
			1. <%= rootPath() %>
			2. <%= userPath({user_id: 1}) %>
			3. <%= myPeepsPath() %>
			5. <%= carPath({car_id: 1}) %>
			6. <%= newCarPath() %>
			7. <%= editCarPath({car_id: 1}) %>
			8. <%= editCarPath({car_id: 1, other: 12}) %>
			9. <%= rootPath({"some":"variable","other": 12}) %>
			10. <%= rootPath() %>
			11. <%= rootPath({"special/":"12=ss"}) %>
			12. <%= resourcePath({resource_id: 1}) %>
			13. <%= editResourcePath({resource_id: 1}) %>
		`))
	}

	a.GET("/", sampleHandler)
	a.GET("/users", sampleHandler)
	a.GET("/users/{user_id}", sampleHandler)
	a.GET("/peeps", sampleHandler).Name("myPeeps")
	a.Resource("/car", carsResource)
	a.Resource("/resources", resourcesResource)

	w := httptest.New(a)
	res := w.HTML("/").Get()

	r.Equal(http.StatusOK, res.Code)
	r.Contains(res.Body.String(), "1. /")
	r.Contains(res.Body.String(), "2. /users/1")
	r.Contains(res.Body.String(), "3. /peeps")
	r.Contains(res.Body.String(), "5. /car/1")
	r.Contains(res.Body.String(), "6. /car/new")
	r.Contains(res.Body.String(), "7. /car/1/edit")
	r.Contains(res.Body.String(), "8. /car/1/edit/?other=12")
	r.Contains(res.Body.String(), "9. /?other=12&some=variable")
	r.Contains(res.Body.String(), "10. /")
	r.Contains(res.Body.String(), "11. /?special%2F=12%3Dss")
	r.Contains(res.Body.String(), "12. /resources/1")
	r.Contains(res.Body.String(), "13. /resources/1/edit")
}

func Test_App_NamedRoutes_MissingParameter(t *testing.T) {
	r := require.New(t)
	a := New(Options{})

	rr := render.New(render.Options{
		HTMLLayout:   "application.plush.html",
		TemplatesBox: packr.New("../templates", "../templates"),
		Helpers:      map[string]interface{}{},
	})

	sampleHandler := func(c Context) error {
		c.Set("opts", map[string]interface{}{})
		return c.Render(http.StatusOK, rr.String(`
			<%= userPath(opts) %>
		`))
	}

	a.GET("/users/{user_id}", sampleHandler)
	w := httptest.New(a)
	res := w.HTML("/users/1").Get()

	r.Equal(http.StatusInternalServerError, res.Code)
	r.Contains(res.Body.String(), "missing parameters for /users/{user_id}")
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

type paramKeyResource struct {
	*userResource
}

func (paramKeyResource) ParamKey() string {
	return "bazKey"
}

func Test_Resource_ParamKey(t *testing.T) {
	r := require.New(t)
	fr := &paramKeyResource{&userResource{}}
	a := New(Options{})
	a.Resource("/foo", fr)
	rt := a.Routes()
	paths := []string{}
	for _, rr := range rt {
		paths = append(paths, rr.Path)
	}
	r.Contains(paths, "/foo/{bazKey}/edit/")
}

type mwResource struct {
	WebResource
}

func (mwResource) Use() []MiddlewareFunc {
	var mw []MiddlewareFunc

	mw = append(mw, func(next Handler) Handler {
		return func(c Context) error {
			if c.Param("good") == "" {
				return fmt.Errorf("not good")
			}
			return next(c)
		}
	})

	return mw
}

func (m mwResource) List(c Context) error {
	return c.Render(http.StatusOK, render.String("southern harmony and the musical companion"))
}

func Test_Resource_MW(t *testing.T) {
	r := require.New(t)
	fr := mwResource{}
	a := New(Options{})
	a.Resource("/foo", fr)

	w := httptest.New(a)
	res := w.HTML("/foo?good=true").Get()
	r.Equal(http.StatusOK, res.Code)
	r.Contains(res.Body.String(), "southern harmony")

	res = w.HTML("/foo").Get()
	r.Equal(http.StatusInternalServerError, res.Code)

	r.NotContains(res.Body.String(), "southern harmony")
}

type userResource struct{}

func (u *userResource) List(c Context) error {
	return c.Render(http.StatusOK, render.String("list"))
}

func (u *userResource) Show(c Context) error {
	return c.Render(http.StatusOK, render.String(`show <%=params["user_id"] %>`))
}

func (u *userResource) New(c Context) error {
	return c.Render(http.StatusOK, render.String("new"))
}

func (u *userResource) Create(c Context) error {
	return c.Render(http.StatusOK, render.String("create"))
}

func (u *userResource) Edit(c Context) error {
	return c.Render(http.StatusOK, render.String(`edit <%=params["user_id"] %>`))
}

func (u *userResource) Update(c Context) error {
	return c.Render(http.StatusOK, render.String(`update <%=params["user_id"] %>`))
}

func (u *userResource) Destroy(c Context) error {
	return c.Render(http.StatusOK, render.String(`destroy <%=params["user_id"] %>`))
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
		"/":                                    "root",
		"/users":                               "users",
		"/users/new":                           "newUsers",
		"/users/{user_id}":                     "user",
		"/users/{user_id}/children":            "userChildren",
		"/users/{user_id}/children/{child_id}": "userChild",
		"/users/{user_id}/children/new":        "newUserChildren",
		"/users/{user_id}/children/{child_id}/build": "userChildBuild",
		"/admin/planes":                 "adminPlanes",
		"/admin/planes/{plane_id}":      "adminPlane",
		"/admin/planes/{plane_id}/edit": "editAdminPlane",
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
		return c.Render(http.StatusOK, rr.String(name))
	})

	w := httptest.New(a)
	res := w.HTML("/john").Get()

	r.Contains(res.Body.String(), "john")
}

func Test_Router_Matches_Trailing_Slash(t *testing.T) {
	table := []struct {
		mapped   string
		browser  string
		expected string
	}{
		{"/foo", "/foo", "/foo/"},
		{"/foo", "/foo/", "/foo/"},
		{"/foo/", "/foo", "/foo/"},
		{"/foo/", "/foo/", "/foo/"},
		{"/index.html", "/index.html", "/index.html/"},
		{"/foo.gif", "/foo.gif", "/foo.gif/"},
		{"/{img}", "/foo.png", "/foo.png/"},
	}

	for _, tt := range table {
		t.Run(tt.mapped+"|"+tt.browser, func(st *testing.T) {
			r := require.New(st)

			app := New(Options{
				PreWares: []PreWare{
					func(h http.Handler) http.Handler {
						var f http.HandlerFunc = func(res http.ResponseWriter, req *http.Request) {
							path := req.URL.Path
							req.URL.Path = strings.TrimSuffix(path, "/")
							r.False(strings.HasSuffix(req.URL.Path, "/"))
							h.ServeHTTP(res, req)
						}
						return f
					},
				},
			})
			app.GET(tt.mapped, func(c Context) error {
				return c.Render(http.StatusOK, render.String(c.Request().URL.Path))
			})

			w := httptest.New(app)
			res := w.HTML(tt.browser).Get()

			r.Equal(http.StatusOK, res.Code)
			r.Equal(tt.expected, res.Body.String())
		})
	}
}
