package buffalo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/stretchr/testify/require"
)

func testApp() http.Handler {
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
