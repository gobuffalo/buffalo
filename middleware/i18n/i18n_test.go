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

type User struct {
	FirstName string
	LastName  string
}

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
	app.GET("/plural", func(c buffalo.Context) error {
		return c.Render(200, r.HTML("plural.html"))
	})
	app.GET("/format", func(c buffalo.Context) error {
		usersList := make([]User, 0)
		usersList = append(usersList, User{"Mark", "Bates"})
		usersList = append(usersList, User{"Chuck", "Berry"})
		c.Set("Users", usersList)
		return c.Render(200, r.HTML("format.html"))
	})
	return app
}

func Test_i18n(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/").Get()
	r.Equal("Hello, World!\n", res.Body.String())
}

func Test_i18n_fr(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	req := w.Request("/")
	// Set language as "french"
	req.Headers["Accept-Language"] = "fr-fr"
	res := req.Get()
	r.Equal("Bonjour à tous !\n", res.Body.String())
}

func Test_i18n_plural(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/plural").Get()
	r.Equal("Hello, alone!\nHello, 5 people!\n", res.Body.String())
}

func Test_i18n_plural_fr(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	req := w.Request("/plural")
	// Set language as "french"
	req.Headers["Accept-Language"] = "fr-fr"
	res := req.Get()
	r.Equal("Bonjour, tout seul !\nBonjour, 5 personnes !\n", res.Body.String())
}

func Test_i18n_format(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/format").Get()
	r.Equal("Hello Mark!\n\n\t* Mr. Mark Bates\n\n\t* Mr. Chuck Berry\n", res.Body.String())
}

func Test_i18n_format_fr(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	req := w.Request("/format")
	// Set language as "french"
	req.Headers["Accept-Language"] = "fr-fr"
	res := req.Get()
	r.Equal("Bonjour Mark !\n\n\t* M. Mark Bates\n\n\t* M. Chuck Berry\n", res.Body.String())
}
