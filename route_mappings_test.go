package buffalo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_App_Routes_without_Root(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	r.Nil(a.root)

	a.GET("/foo", voidHandler)

	routes := a.Routes()
	r.Len(routes, 1)
	route := routes[0]
	r.Equal("GET", route.Method)
	r.Equal("/foo/", route.Path)
	r.NotZero(route.HandlerName)
}

type resourceHandler struct{}

func (r resourceHandler) List(Context) error {
	return nil
}

func (r resourceHandler) Show(Context) error {
	return nil
}

func (r resourceHandler) Create(Context) error {
	return nil
}

func (r resourceHandler) Update(Context) error {
	return nil
}

func (r resourceHandler) Destroy(Context) error {
	return nil
}

func Test_App_Routes_Resource(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	r.Nil(a.root)

	a.GET("/foo", voidHandler)
	a.Resource("/r", resourceHandler{})

	routes := a.Routes()
	r.Len(routes, 6)
	route := routes[0]
	r.Equal("GET", route.Method)
	r.Equal("/foo/", route.Path)
	r.NotZero(route.HandlerName)

	for k, v := range routes {
		if k > 0 {
			r.Equal("resourceHandler", v.ResourceName)
		}
	}
}

func Test_App_Routes_with_Root(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	r.Nil(a.root)

	g := a.Group("/api/v1")
	g.GET("/foo", voidHandler)

	routes := a.Routes()
	r.Len(routes, 1)
	route := routes[0]
	r.Equal("GET", route.Method)
	r.Equal("/api/v1/foo/", route.Path)
	r.NotZero(route.HandlerName)

	r.Equal(a.Routes(), g.Routes())
}

func Test_App_RouteName(t *testing.T) {
	r := require.New(t)

	a := New(Options{})

	cases := map[string]string{
		"cool":                "coolPath",
		"coolPath":            "coolPath",
		"coco_path":           "cocoPath",
		"ouch_something_cool": "ouchSomethingCoolPath",
	}

	ri := a.GET("/something", voidHandler)
	for k, v := range cases {
		ri.Name(k)
		r.Equal(ri.PathName, v)
	}

}

func Test_RouteList_Lookup(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	r.Nil(a.root)

	a.GET("/foo", voidHandler)
	a.GET("/test", voidHandler)

	routes := a.Routes()
	for _, route := range routes {
		lRoute, err := routes.Lookup(route.PathName)
		r.NoError(err)
		r.Equal(lRoute, route)
	}
	lRoute, err := routes.Lookup("a")
	r.Error(err)
	r.Nil(lRoute)

}

func Test_App_RouteHelpers(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	r.Nil(a.root)

	a.GET("/foo", voidHandler)
	a.GET("/test/{id}", voidHandler)

	rh := a.RouteHelpers()

	r.Len(rh, 2)

	f, ok := rh["fooPath"]
	r.True(ok)
	x, err := f(map[string]interface{}{})
	r.NoError(err)
	r.Equal("/foo/", string(x))

	f, ok = rh["testPath"]
	r.True(ok)
	x, err = f(map[string]interface{}{
		"id": 1,
	})
	r.NoError(err)
	r.Equal("/test/1/", string(x))
}
