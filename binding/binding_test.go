package binding

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Register(t *testing.T) {
	r := require.New(t)
	l := len(defaultRequestBinder.binders)
	Register("foo/bar", func(*http.Request, interface{}) error {
		return nil
	})

	r.Len(defaultRequestBinder.binders, l+1)
}

func Test_RegisterCustomDecoder(t *testing.T) {
	// r := require.New(t)

	// RegisterCustomDecoder(func(vals []string) (interface{}, error) {
	// 	return []string{"X"}, nil
	// }, []interface{}{[]string{}}, nil)

	// type Xt struct {
	// 	Vals []string
	// }

	// type U struct {
	// 	Xt Xt
	// }

	// var ux U
	// app := buffalo.New(buffalo.Options{})
	// app.POST("/", func(c buffalo.Context) error {
	// 	return c.Bind(&ux)
	// })

	// w := httptest.New(app)
	// res := w.HTML("/").Post(&U{
	// 	Xt: Xt{[]string{"foo"}},
	// })

	// r.Equal(http.StatusOK, res.Code)
	// r.Equal([]string{"X"}, ux.Xt.Vals)
}
