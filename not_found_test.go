package buffalo

import (
	"encoding/json"
	"testing"

	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_App_Dev_NotFound(t *testing.T) {
	r := require.New(t)

	a := Automatic(Options{})
	a.Env = "development"
	a.GET("/foo", func(c Context) error { return nil })

	w := willie.New(a)
	res := w.Request("/bad").Get()

	body := res.Body.String()
	r.Contains(body, "404 PAGE NOT FOUND")
	r.Contains(body, "/foo")
	r.Equal(404, res.Code)
}

func Test_App_Dev_NotFound_JSON(t *testing.T) {
	r := require.New(t)

	a := Automatic(Options{})
	a.Env = "development"
	a.GET("/foo", func(c Context) error { return nil })

	w := willie.New(a)
	res := w.JSON("/bad").Get()
	r.Equal(404, res.Code)

	jb := map[string]interface{}{}
	err := json.NewDecoder(res.Body).Decode(&jb)
	r.NoError(err)
	r.Equal("GET", jb["method"])
	r.Equal("/bad", jb["path"])
	r.NotEmpty(jb["routes"])
}

func Test_App_Prod_NotFound(t *testing.T) {
	r := require.New(t)

	a := New(Options{Env: "production"})
	a.GET("/foo", func(c Context) error { return nil })

	w := willie.New(a)
	res := w.Request("/bad").Get()
	r.Equal(404, res.Code)

	body := res.Body.String()
	r.Equal(body, "404 page not found\n")
	r.NotContains(body, "/foo")
}

func Test_App_Override_NotFound(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.ErrorHandlers[404] = func(status int, err error, c Context) error {
		c.Response().WriteHeader(404)
		c.Response().Write([]byte("oops!!!"))
		return nil
	}
	a.GET("/foo", func(c Context) error { return nil })

	w := willie.New(a)
	res := w.Request("/bad").Get()
	r.Equal(404, res.Code)

	body := res.Body.String()
	r.Equal(body, "oops!!!")
	r.NotContains(body, "/foo")
}
