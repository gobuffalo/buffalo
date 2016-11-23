package buffalo

import (
	"net/http"

	"github.com/markbates/buffalo/render"
)

func testApp() http.Handler {
	a := New(Options{})
	rt := a.Group("/router/tests")
	rt.GET("/get", func(c Context) error {
		return c.Render(200, render.String("GET"))
	})
	return a
}

// func Test_Router(t *testing.T) {
// 	r := require.New(t)
//
// 	table := map[string]string{
// 		"get": "GET",
// 	}
//
// 	for k, v := range table {
// 		req := http.NewRequest(v, fmt.Sprintf("/router/tests/%s", k))
// 		res, err := http.DefaultClient.Do(req)
// 		r.NoError(err)
// 	}
// }
