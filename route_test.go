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
	r.Equal("/foo", route.Path)
	r.NotZero(route.HandlerName)
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
	r.Equal("/api/v1/foo", route.Path)
	r.NotZero(route.HandlerName)

	r.Equal(a.Routes(), g.Routes())
}
