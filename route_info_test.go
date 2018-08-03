package buffalo

import (
	"database/sql"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_RouteInfo_ServeHTTP_SQL_Error(t *testing.T) {
	r := require.New(t)

	app := New(Options{})
	app.GET("/good", func(c Context) error {
		return c.Render(200, render.String("hi"))
	})

	app.GET("/bad", func(c Context) error {
		return sql.ErrNoRows
	})

	w := willie.New(app)

	res := w.HTML("/good").Get()
	r.Equal(200, res.Code)

	res = w.HTML("/bad").Get()
	r.Equal(404, res.Code)
}
