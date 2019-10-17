package buffalo

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

func Test_FlashAdd(t *testing.T) {
	r := require.New(t)
	f := newFlash(&Session{})

	r.Equal(f.data, map[string][]string{})

	f.Add("error", "something")
	r.Equal(f.data, map[string][]string{
		"error": {"something"},
	})

	f.Add("error", "other")
	r.Equal(f.data, map[string][]string{
		"error": {"something", "other"},
	})
}

func Test_FlashRender(t *testing.T) {
	r := require.New(t)
	a := New(Options{})
	rr := render.New(render.Options{})

	a.GET("/", func(c Context) error {
		c.Flash().Add("errors", "Error AJ set")
		c.Flash().Add("errors", "Error DAL set")

		return c.Render(http.StatusCreated, rr.String(errorsTPL))
	})

	w := httptest.New(a)
	res := w.HTML("/").Get()

	r.Contains(res.Body.String(), "Error AJ set")
	r.Contains(res.Body.String(), "Error DAL set")
}

func Test_FlashRenderEmpty(t *testing.T) {
	r := require.New(t)
	a := New(Options{})
	rr := render.New(render.Options{})

	a.GET("/", func(c Context) error {
		return c.Render(http.StatusCreated, rr.String(errorsTPL))
	})

	w := httptest.New(a)

	res := w.HTML("/").Get()
	r.NotContains(res.Body.String(), "Flash:")
}

const errorsTPL = `
<%= for (k, v) in flash["errors"] { %>
	Flash:
		<%= k %>:<%= v %>
<% } %>
`

func Test_FlashRenderEntireFlash(t *testing.T) {
	r := require.New(t)
	a := New(Options{})
	rr := render.New(render.Options{})

	a.GET("/", func(c Context) error {
		c.Flash().Add("something", "something to say!")
		return c.Render(http.StatusCreated, rr.String(keyTPL))
	})

	w := httptest.New(a)
	res := w.HTML("/").Get()
	r.Contains(res.Body.String(), "something to say!")
}

const keyTPL = `<%= for (k, v) in flash { %>
	Flash:
		<%= k %>:<%= v %>
<% } %>
`

func Test_FlashRenderCustomKey(t *testing.T) {
	r := require.New(t)
	a := New(Options{})
	rr := render.New(render.Options{})

	a.GET("/", func(c Context) error {
		c.Flash().Add("something", "something to say!")
		return c.Render(http.StatusCreated, rr.String(keyTPL))
	})

	w := httptest.New(a)
	res := w.HTML("/").Get()
	r.Contains(res.Body.String(), "something to say!")
}

func Test_FlashRenderCustomKeyNotDefined(t *testing.T) {
	r := require.New(t)
	a := New(Options{})
	rr := render.New(render.Options{})

	a.GET("/", func(c Context) error {
		return c.Render(http.StatusCreated, rr.String(customKeyTPL))
	})

	w := httptest.New(a)
	res := w.HTML("/").Get()
	r.NotContains(res.Body.String(), "something to say!")
}

const customKeyTPL = `
	{{#each flash.other as |k value|}}
		{{value}}
	{{/each}}`
