package buffalo

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func TestHandleRobotsNoFile(t *testing.T) {
	r := require.New(t)
	rend := render.New(render.Options{
		AssetsBox: packr.NewBox(os.TempDir()),
	})

	app := New(Options{})
	app.GET("/robots.txt", NewRobotsHandler(rend))

	w := willie.New(app)
	res := w.Request("/robots.txt").Get()

	r.Contains(res.Body.String(), "User-agent: *\nDisallow:")
}

func TestHandleRobotsWithFile(t *testing.T) {
	r := require.New(t)
	dir := os.TempDir()
	d1 := []byte("User-agent: *\nDisallow: /cgi-bin/\nDisallow: /tmp/\nDisallow: /junk/")

	err := ioutil.WriteFile(filepath.Join(dir, "robots.txt"), d1, 0644)
	if err != nil {
		r.Fail("Could not create robots file")
	}

	rend := render.New(render.Options{
		AssetsBox: packr.NewBox(dir),
	})

	app := New(Options{})
	app.GET("/robots.txt", NewRobotsHandler(rend))

	w := willie.New(app)
	res := w.Request("/robots.txt").Get()

	r.Equal(res.Body.String(), "User-agent: *\nDisallow: /cgi-bin/\nDisallow: /tmp/\nDisallow: /junk/")
}

func TestHandleRobotsWithNoBox(t *testing.T) {
	r := require.New(t)
	rend := render.New(render.Options{})
	app := New(Options{})
	app.GET("/robots.txt", NewRobotsHandler(rend))

	w := willie.New(app)
	res := w.Request("/robots.txt").Get()

	r.Equal(res.Body.String(), "User-agent: *\nDisallow:")
}
