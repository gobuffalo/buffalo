package i18n

import (
	"log"
	"strings"
	"testing"
	"time"

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
	// Setup URL prefix Language extractor
	t.LanguageExtractors = append(t.LanguageExtractors, URLPrefixLanguageExtractor)

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
	app.GET("/collision", func(c buffalo.Context) error {
		return c.Render(200, r.HTML("collision.html"))
	})
	app.GET("/localized", func(c buffalo.Context) error {
		return c.Render(200, r.HTML("localized_view.html"))
	})
	app.GET("/languages-list", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(t.AvailableLanguages()))
	})
	app.GET("/refresh", func(c buffalo.Context) error {
		// This flash will be displayed in english
		c.Flash().Add("success", t.Translate(c, "refresh-success"))

		// Change lang to fr-fr
		c.Cookies().Set("lang", "fr-fr", time.Minute)
		t.Refresh(c, "fr-fr")

		// This flash will be displayed in french
		c.Flash().Add("success", t.Translate(c, "refresh-success"))
		return c.Render(200, r.HTML("refresh.html"))
	})
	// Disable i18n middleware
	noI18n := func(c buffalo.Context) error {
		return c.Render(200, r.HTML("localized_view.html"))
	}
	app.Middleware.Skip(t.Middleware(), noI18n)
	app.GET("/localized-disabled", noI18n)
	app.GET("/{lang:fr|en}/index", func(c buffalo.Context) error {
		return c.Render(200, r.HTML("index.html"))
	})
	return app
}

func Test_i18n(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/").Get()
	r.Equal("Hello, World!", strings.TrimSpace(res.Body.String()))
}

func Test_i18n_fr(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	req := w.Request("/")
	// Set language as "french"
	req.Headers["Accept-Language"] = "fr-fr"
	res := req.Get()
	r.Equal("Bonjour à tous !", strings.TrimSpace(res.Body.String()))
}

func Test_i18n_plural(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/plural").Get()
	r.Equal("Hello, alone!\nHello, 5 people!", strings.TrimSpace(res.Body.String()))
}

func Test_i18n_plural_fr(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	req := w.Request("/plural")
	// Set language as "french"
	req.Headers["Accept-Language"] = "fr-fr"
	res := req.Get()
	r.Equal("Bonjour, tout seul !\nBonjour, 5 personnes !", strings.TrimSpace(res.Body.String()))
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

func Test_i18n_Localized_View(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	// Test with complex Accept-Language
	req := w.Request("/localized")
	req.Headers["Accept-Language"] = "en-UK,en-US;q=0.5"
	res := req.Get()
	r.Equal("Hello!", strings.TrimSpace(res.Body.String()))

	// Test priority
	req.Headers["Accept-Language"] = "fr,en-US"
	res = req.Get()
	r.Equal("Bonjour !", strings.TrimSpace(res.Body.String()))

	// Test fallback
	req.Headers["Accept-Language"] = "ru"
	res = req.Get()
	r.Equal("Default", strings.TrimSpace(res.Body.String()))

	// Test i18n disabled
	req = w.Request("/localized-disabled")
	req.Headers["Accept-Language"] = "en-UK,en-US;q=0.5"
	res = req.Get()
	r.Equal("Default", strings.TrimSpace(res.Body.String()))
}

func Test_i18n_collision(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/collision").Get()
	r.Equal("Collision OK", strings.TrimSpace(res.Body.String()))
}

func Test_i18n_availableLanguages(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	res := w.Request("/languages-list").Get()
	r.Equal("[\"en-us\",\"fr-fr\"]", strings.TrimSpace(res.Body.String()))
}

func Test_i18n_URL_prefix(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	req := w.Request("/fr/index")
	res := req.Get()
	r.Equal("Bonjour à tous !", strings.TrimSpace(res.Body.String()))

	req = w.Request("/en/index")
	res = req.Get()
	r.Equal("Hello, World!", strings.TrimSpace(res.Body.String()))
}

func Test_Refresh(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())
	req := w.Request("/refresh")
	res := req.Get()
	r.Equal("success: Language changed!#success: Langue modifiée !#", strings.TrimSpace(res.Body.String()))
}
