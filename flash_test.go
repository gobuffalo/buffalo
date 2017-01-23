package buffalo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_FlashAdd(t *testing.T) {
	r := require.New(t)
	f := newFlash(&Session{})

	r.Equal(f.data, map[string][]string{})

	f.Add("error", "something")
	r.Equal(f.data, map[string][]string{
		"error": []string{"something"},
	})

	f.Add("error", "other")
	r.Equal(f.data, map[string][]string{
		"error": []string{"something", "other"},
	})
}

func Test_FlashRender(t *testing.T) {
	tempFolder := os.TempDir()
	ioutil.WriteFile(tempFolder+"application.html", []byte("{{yield}}"), 0755)
	ioutil.WriteFile(tempFolder+"show.html", []byte(errorsTPL), 0755)

	r := require.New(t)
	a := Automatic(Options{})
	rr := render.New(render.Options{
		HTMLLayout: tempFolder + "application.html",
	})

	a.GET("/", func(c Context) error {
		c.Flash().Add("errors", "Error AJ set")
		c.Flash().Add("errors", "Error DAL set")

		return c.Render(201, rr.HTML(tempFolder+"show.html"))
	})

	w := willie.New(a)
	res := w.Request("/").Get()

	r.Contains(res.Body.String(), "Error AJ set")
	r.Contains(res.Body.String(), "Error DAL set")
}

func Test_FlashRenderEmpty(t *testing.T) {
	tempFolder := os.TempDir()
	ioutil.WriteFile(tempFolder+"application.html", []byte("{{yield}}"), 0755)
	ioutil.WriteFile(tempFolder+"show.html", []byte(errorsTPL), 0755)

	r := require.New(t)
	a := Automatic(Options{})
	rr := render.New(render.Options{
		HTMLLayout: tempFolder + "application.html",
	})

	a.GET("/", func(c Context) error {
		return c.Render(201, rr.HTML(tempFolder+"show.html"))
	})

	w := willie.New(a)

	res := w.Request("/").Get()
	r.NotContains(res.Body.String(), "Flash:")
}

const errorsTPL = `{{#each flash.errors as |k flash|}}
	Flash:
    {{k}}:{{flash}}
{{/each}}`

func Test_FlashRenderEntireFlash(t *testing.T) {
	tempFolder := os.TempDir()
	ioutil.WriteFile(tempFolder+"application.html", []byte("{{yield}}"), 0755)
	ioutil.WriteFile(tempFolder+"show.html", []byte(keyTPL), 0755)

	r := require.New(t)
	a := Automatic(Options{})
	rr := render.New(render.Options{
		HTMLLayout: tempFolder + "application.html",
	})

	a.GET("/", func(c Context) error {
		c.Flash().Add("something", "something to say!")
		return c.Render(201, rr.HTML(tempFolder+"show.html"))
	})

	w := willie.New(a)
	res := w.Request("/").Get()
	r.Contains(res.Body.String(), "something to say!")
}

const keyTPL = `{{#each flash as |k flash|}}
	Flash:
    {{k}}:{{flash}}
{{/each}}`
