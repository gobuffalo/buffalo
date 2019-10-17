package buffalo

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

func Test_App_Dev_NotFound(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.Env = "development"
	a.GET("/foo", func(c Context) error { return nil })

	w := httptest.New(a)
	res := w.HTML("/bad").Get()

	body := res.Body.String()
	r.Contains(body, "404 - ERROR!")
	r.Contains(body, "/foo")
	r.Equal(http.StatusNotFound, res.Code)
}

func Test_App_Dev_NotFound_JSON(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.Env = "development"
	a.GET("/foo", func(c Context) error { return nil })

	w := httptest.New(a)
	res := w.JSON("/bad").Get()
	r.Equal(http.StatusNotFound, res.Code)

	jb := map[string]interface{}{}
	err := json.NewDecoder(res.Body).Decode(&jb)
	r.NoError(err)
	r.Equal(float64(http.StatusNotFound), jb["code"])
}

func Test_App_Override_NotFound(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	a.ErrorHandlers[http.StatusNotFound] = func(status int, err error, c Context) error {
		c.Response().WriteHeader(http.StatusNotFound)
		c.Response().Write([]byte("oops!!!"))
		return nil
	}
	a.GET("/foo", func(c Context) error { return nil })

	w := httptest.New(a)
	res := w.HTML("/bad").Get()
	r.Equal(http.StatusNotFound, res.Code)

	body := res.Body.String()
	r.Equal(body, "oops!!!")
	r.NotContains(body, "/foo")
}
