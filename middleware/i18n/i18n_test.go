package i18n

import (
	"log"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func app() *buffalo.App {
	app := buffalo.New(buffalo.Options{})

	r := render.New(render.Options{
		TemplatesBox: packr.NewBox("./templates"),
	})

	// Setup and use translations:
	t, err := New(packr.NewBox("./locales"), "en-US")
	if err != nil {
		log.Fatal(err)
	}
	app.Use(t.Middleware())
	app.GET("/", func(c buffalo.Context) error {
		return c.Render(200, r.HTML("index.html"))
	})
	return app
}

func Test_i18n(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/").Get()
	r.Equal("Hello, World!\n", res.Body.String())
}
