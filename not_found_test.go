package buffalo

import (
	"net/http"
	"testing"

	"github.com/markbates/going/willy"
	"github.com/stretchr/testify/require"
)

func Test_App_Dev_NotFound(t *testing.T) {
	r := require.New(t)

	a := New(Options{Env: "development"})
	a.GET("/foo", func(c Context) error { return nil })

	w := willy.New(a)
	res := w.Request("/bad").Get()
	r.Equal(404, res.Code)

	body := res.Body.String()
	r.Contains(body, "404 PAGE NOT FOUND")
	r.Contains(body, "/foo")
}

func Test_App_Prod_NotFound(t *testing.T) {
	r := require.New(t)

	a := New(Options{Env: "production"})
	a.GET("/foo", func(c Context) error { return nil })

	w := willy.New(a)
	res := w.Request("/bad").Get()
	r.Equal(404, res.Code)

	body := res.Body.String()
	r.Equal(body, "404 page not found\n")
	r.NotContains(body, "/foo")
}

func Test_App_Override_NotFound(t *testing.T) {
	r := require.New(t)

	a := New(Options{
		NotFound: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(404)
			res.Write([]byte("oops!!!"))
		}),
	})
	a.GET("/foo", func(c Context) error { return nil })

	w := willy.New(a)
	res := w.Request("/bad").Get()
	r.Equal(404, res.Code)

	body := res.Body.String()
	r.Equal(body, "oops!!!")
	r.NotContains(body, "/foo")
}
