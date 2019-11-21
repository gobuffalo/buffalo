package binding_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

func Test_BindTwiceJSON(t *testing.T) {
	r := require.New(t)

	render := render.New(render.Options{
		DefaultContentType: "application/json",
	})

	type U struct {
		F string
	}
	app := buffalo.New(buffalo.Options{})
	app.POST("/", func(c buffalo.Context) error {
		var ux, uy U
		// Bind once
		if err := c.Bind(&ux); err != nil {
			return c.Render(http.StatusInternalServerError, render.String(err.Error()))
		}
		// Bind twice
		if err := c.Bind(&uy); err != nil {
			return c.Render(http.StatusInternalServerError, render.String(err.Error()))
		}
		return c.Render(http.StatusOK, render.JSON(U{
			F: "ux=" + ux.F + ", uy=" + uy.F + " OK",
		}))
	})

	w := httptest.New(app)
	res := w.JSON("/").Post(&U{F: "foo"})

	r.Equal(http.StatusOK, res.Code, "Http code not OK, body='%v'", res.Body)

	var resVal U
	r.NoError(json.Unmarshal(res.Body.Bytes(), &resVal))
	r.Equal("ux=foo, uy=foo OK", resVal.F)
}
