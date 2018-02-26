package render_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

type Car struct {
	ID   int
	Name string
}

type Cars []Car

func Test_Auto_JSON(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars", func(c buffalo.Context) error {
		return c.Render(200, render.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
	})

	w := willie.New(app)

	res := w.JSON("/cars").Get()
	r.Equal(`["Honda","Toyota","Ford","Chevy"]`, strings.TrimSpace(res.Body.String()))
}

func Test_Auto_XML(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars", func(c buffalo.Context) error {
		return c.Render(200, render.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
	})

	w := willie.New(app)

	res := w.XML("/cars").Get()
	r.Equal("<string>Honda</string>\n<string>Toyota</string>\n<string>Ford</string>\n<string>Chevy</string>", strings.TrimSpace(res.Body.String()))
}

func Test_Auto_HTML_List(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/index.html", p), "INDEX: <%= len(cars) %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.GET(fmt.Sprintf("/%s", p), func(c buffalo.Context) error {
				return c.Render(200, e.Auto(c, Cars{
					{Name: "Ford"},
					{Name: "Chevy"},
				}))
			})

			w := willie.New(app)
			res := w.HTML("/%s", p).Get()

			r.Contains(res.Body.String(), "INDEX: 2")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_Show(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/show.html", p), "Show: <%= car.Name %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.GET(fmt.Sprintf("/%s/{id}", p), func(c buffalo.Context) error {
				return c.Render(200, e.Auto(c, Car{Name: "Honda"}))
			})

			w := willie.New(app)
			res := w.HTML("/%s/1", p).Get()
			r.Contains(res.Body.String(), "Show: Honda")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_New(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/new.html", p), "New: <%= car.Name %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.GET(fmt.Sprintf("/%s/new", p), func(c buffalo.Context) error {
				return c.Render(200, e.Auto(c, Car{Name: "Honda"}))
			})

			w := willie.New(app)
			res := w.HTML("/%s/new", p).Get()
			r.Contains(res.Body.String(), "New: Honda")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_Create(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/new.html", p), "New: <%= car.Name %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.POST(fmt.Sprintf("/%s", p), func(c buffalo.Context) error {
				return c.Render(201, e.Auto(c, Car{Name: "Honda"}))
			})

			w := willie.New(app)
			res := w.HTML("/%s", p).Post(nil)
			r.Contains(res.Body.String(), "New: Honda")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_Create_Redirect(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		app := buffalo.New(buffalo.Options{})
		app.POST(fmt.Sprintf("/%s", p), func(c buffalo.Context) error {
			return c.Render(201, render.Auto(c, Car{
				ID:   1,
				Name: "Honda",
			}))
		})

		w := willie.New(app)
		res := w.HTML("/%s", p).Post(nil)
		r.Equal(fmt.Sprintf("/%s/1", p), res.Location())
		r.Equal(302, res.Code)
	}
}

func Test_Auto_HTML_Create_Redirect_Error(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/new.html", p), "Create: <%= car.Name %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.POST(fmt.Sprintf("/%s", p), func(c buffalo.Context) error {
				b := Car{
					Name: "Honda",
				}
				return c.Render(422, e.Auto(c, b))
			})

			w := willie.New(app)
			res := w.HTML("/%s", p).Post(nil)
			r.Equal(422, res.Code)
			r.Contains(res.Body.String(), "Create: Honda")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_Edit(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/edit.html", p), "Edit: <%= car.Name %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.GET(fmt.Sprintf("/%s/{id}/edit", p), func(c buffalo.Context) error {
				return c.Render(200, e.Auto(c, Car{Name: "Honda"}))
			})

			w := willie.New(app)
			res := w.HTML("/%s/1/edit", p).Get()
			r.Contains(res.Body.String(), "Edit: Honda")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_Update(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/edit.html", p), "Update: <%= car.Name %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.PUT(fmt.Sprintf("/%s/{id}", p), func(c buffalo.Context) error {
				return c.Render(422, e.Auto(c, Car{Name: "Honda"}))
			})

			w := willie.New(app)
			res := w.HTML("/%s/1", p).Put(nil)

			r.Contains(res.Body.String(), "Update: Honda")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_Update_Redirect(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		app := buffalo.New(buffalo.Options{})
		app.PUT(fmt.Sprintf("/%s/{id}", p), func(c buffalo.Context) error {
			b := Car{
				ID:   1,
				Name: "Honda",
			}
			return c.Render(200, render.Auto(c, b))
		})

		w := willie.New(app)
		res := w.HTML("/%s/1", p).Put(nil)
		r.Equal(fmt.Sprintf("/%s/1", p), res.Location())
		r.Equal(302, res.Code)
	}
}

func Test_Auto_HTML_Update_Redirect_Error(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		err := withHTMLFile(fmt.Sprintf("%s/edit.html", p), "Update: <%= car.Name %>", func(e *render.Engine) {
			app := buffalo.New(buffalo.Options{})
			app.PUT(fmt.Sprintf("/%s/{id}", p), func(c buffalo.Context) error {
				b := Car{
					ID:   1,
					Name: "Honda",
				}
				return c.Render(422, e.Auto(c, b))
			})

			w := willie.New(app)
			res := w.HTML("/%s/1", p).Put(nil)
			r.Equal(422, res.Code)
			r.Contains(res.Body.String(), "Update: Honda")
		})
		r.NoError(err)
	}
}

func Test_Auto_HTML_Destroy_Redirect(t *testing.T) {
	r := require.New(t)

	for _, p := range []string{"cars", "admin/cars"} {
		app := buffalo.New(buffalo.Options{})
		app.DELETE(fmt.Sprintf("/%s/{id}", p), func(c buffalo.Context) error {
			b := Car{
				ID:   1,
				Name: "Honda",
			}
			return c.Render(200, render.Auto(c, b))
		})

		w := willie.New(app)
		res := w.HTML("/%s/1", p).Delete()
		r.Equal(fmt.Sprintf("/%s", p), res.Location())
		r.Equal(302, res.Code)
	}
}

func withHTMLFile(name string, contents string, fn func(*render.Engine)) error {
	tmpDir := filepath.Join(os.TempDir(), filepath.Dir(name))
	err := os.MkdirAll(tmpDir, 0766)
	if err != nil {
		return err
	}
	defer os.Remove(tmpDir)

	tmpFile, err := os.Create(filepath.Join(tmpDir, filepath.Base(name)))
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(contents))
	if err != nil {
		return err
	}

	e := render.New(render.Options{
		TemplatesBox: packr.NewBox(os.TempDir()),
	})

	fn(e)
	return nil
}
