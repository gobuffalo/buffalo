package binding_test

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

type Xt struct {
	Vals []string
}

func Test_RegisterCustomDecoder(t *testing.T) {
	r := require.New(t)

	binding.RegisterCustomDecoder(func(vals []string) (interface{}, error) {
		return []string{"X"}, nil
	}, []interface{}{[]string{}}, nil)

	type U struct {
		Xt Xt
	}
	var ux U
	app := buffalo.New(buffalo.Options{})
	app.POST("/", func(c buffalo.Context) error {
		return c.Bind(&ux)
	})

	w := httptest.New(app)
	res := w.HTML("/").Post(&U{
		Xt: Xt{[]string{"foo"}},
	})
	r.Equal(http.StatusOK, res.Code)

	r.Equal([]string{"X"}, ux.Xt.Vals)
}
