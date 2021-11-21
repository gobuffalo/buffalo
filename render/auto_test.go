package render_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/psanford/memfs"
	"github.com/stretchr/testify/require"
)

type Car struct {
	ID   int
	Name string
}

type Cars []Car

func Test_Auto_DefaultContentType(t *testing.T) {
	r := require.New(t)

	re := render.New(render.Options{
		DefaultContentType: "application/json",
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
	})

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/cars", nil)
	app.ServeHTTP(res, req)

	r.Equal(`["Honda","Toyota","Ford","Chevy"]`, strings.TrimSpace(res.Body.String()))
}

func Test_Auto_JSON(t *testing.T) {
	r := require.New(t)

	re := render.New(render.Options{})
	app := buffalo.New(buffalo.Options{})
	app.GET("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
	})

	w := httptest.New(app)

	res := w.JSON("/cars").Get()
	r.Equal(`["Honda","Toyota","Ford","Chevy"]`, strings.TrimSpace(res.Body.String()))
}

func Test_Auto_XML(t *testing.T) {
	r := require.New(t)

	re := render.New(render.Options{})
	app := buffalo.New(buffalo.Options{})
	app.GET("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
	})

	w := httptest.New(app)

	res := w.XML("/cars").Get()
	r.Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<string>Honda</string>\n<string>Toyota</string>\n<string>Ford</string>\n<string>Chevy</string>", strings.TrimSpace(res.Body.String()))
}

func Test_Auto_HTML_List(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/index.html", []byte("INDEX: <%= len(cars) %>"), 0644))

	re := render.NewEngine()
	re.TemplatesFS = rootFS

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, Cars{
			{Name: "Ford"},
			{Name: "Chevy"},
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars").Get()

	r.Contains(res.Body.String(), "INDEX: 2")
}

func Test_Auto_HTML_List_Plural(t *testing.T) {
	r := require.New(t)

	type Person struct {
		Name string
	}

	type People []Person

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("people", 0755))
	r.NoError(rootFS.WriteFile("people/index.html", []byte("INDEX: <%= len(people) %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/people", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, People{
			Person{Name: "Ford"},
			Person{Name: "Chevy"},
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/people").Get()

	r.Contains(res.Body.String(), "INDEX: 2")
}

func Test_Auto_HTML_Show(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/show.html", []byte("Show: <%= car.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars/{id}", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Get()
	r.Contains(res.Body.String(), "Show: Honda")
}

func Test_Auto_HTML_New(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/new.html", []byte("New: <%= car.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars/new", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/new").Get()
	r.Contains(res.Body.String(), "New: Honda")
}

func Test_Auto_HTML_Create(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/new.html", []byte("New: <%= car.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.POST("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars").Post(nil)
	r.Contains(res.Body.String(), "New: Honda")
}

func Test_Auto_HTML_Create_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.POST("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, render.Auto(c, Car{
			ID:   1,
			Name: "Honda",
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars").Post(nil)
	r.Equal("/cars/1", res.Location())
	r.Equal(http.StatusFound, res.Code)
}

func Test_Auto_HTML_Create_Redirect_Error(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/new.html", []byte("Create: <%= car.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.POST("/cars", func(c buffalo.Context) error {
		b := Car{
			Name: "Honda",
		}
		return c.Render(http.StatusUnprocessableEntity, re.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars").Post(nil)
	r.Equal(http.StatusUnprocessableEntity, res.Code)
	r.Contains(res.Body.String(), "Create: Honda")
}

func Test_Auto_HTML_Create_Nested_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	admin := app.Group("/admin")
	admin.POST("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, render.Auto(c, Car{
			ID:   1,
			Name: "Honda",
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/admin/cars").Post(nil)
	r.Equal("/admin/cars/1", res.Location())
	r.Equal(http.StatusFound, res.Code)
}

func Test_Auto_HTML_Destroy_Nested_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	admin := app.Group("/admin")
	admin.DELETE("/cars", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, render.Auto(c, Car{
			ID:   1,
			Name: "Honda",
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/admin/cars").Delete()
	r.Equal("/admin/cars", res.Location())
	r.Equal(http.StatusFound, res.Code)
}

func Test_Auto_HTML_Edit(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/edit.html", []byte("Edit: <%= car.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars/{id}/edit", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1/edit").Get()
	r.Contains(res.Body.String(), "Edit: Honda")
}

func Test_Auto_HTML_Update(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/edit.html", []byte("Update: <%= car.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.PUT("/cars/{id}", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Put(nil)

	r.Contains(res.Body.String(), "Update: Honda")
}

func Test_Auto_HTML_Update_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.PUT("/cars/{id}", func(c buffalo.Context) error {
		b := Car{
			ID:   1,
			Name: "Honda",
		}
		return c.Render(http.StatusOK, render.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Put(nil)
	r.Equal("/cars/1", res.Location())
	r.Equal(http.StatusFound, res.Code)
}

func Test_Auto_HTML_Update_Redirect_Error(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("cars", 0755))
	r.NoError(rootFS.WriteFile("cars/edit.html", []byte("Update: <%= car.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.PUT("/cars/{id}", func(c buffalo.Context) error {
		b := Car{
			ID:   1,
			Name: "Honda",
		}
		return c.Render(http.StatusUnprocessableEntity, re.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Put(nil)
	r.Equal(http.StatusUnprocessableEntity, res.Code)
	r.Contains(res.Body.String(), "Update: Honda")
}

func Test_Auto_HTML_Destroy_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.DELETE("/cars/{id}", func(c buffalo.Context) error {
		b := Car{
			ID:   1,
			Name: "Honda",
		}
		return c.Render(http.StatusOK, render.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Delete()
	r.Equal("/cars/", res.Location())
	r.Equal(http.StatusFound, res.Code)
}

func Test_Auto_HTML_List_Plural_MultiWord(t *testing.T) {
	r := require.New(t)

	type RoomProvider struct {
		Name string
	}

	type RoomProviders []RoomProvider

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("room_providers", 0755))
	r.NoError(rootFS.WriteFile("room_providers/index.html", []byte("INDEX: <%= len(roomProviders) %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/room_providers", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, RoomProviders{
			RoomProvider{Name: "Ford"},
			RoomProvider{Name: "Chevy"},
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/room_providers").Get()

	r.Contains(res.Body.String(), "INDEX: 2")
}

func Test_Auto_HTML_List_Plural_MultiWord_Dashed(t *testing.T) {
	r := require.New(t)

	type RoomProvider struct {
		Name string
	}

	type RoomProviders []RoomProvider

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("room_providers", 0755))
	r.NoError(rootFS.WriteFile("room_providers/index.html", []byte("INDEX: <%= len(roomProviders) %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/room-providers", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, RoomProviders{
			RoomProvider{Name: "Ford"},
			RoomProvider{Name: "Chevy"},
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/room-providers").Get()

	r.Contains(res.Body.String(), "INDEX: 2")
}

func Test_Auto_HTML_Show_MultiWord_Dashed(t *testing.T) {
	r := require.New(t)

	type RoomProvider struct {
		ID   int
		Name string
	}

	rootFS := memfs.New()
	r.NoError(rootFS.MkdirAll("room_providers", 0755))
	r.NoError(rootFS.WriteFile("room_providers/show.html", []byte("SHOW: <%= roomProvider.Name %>"), 0644))

	re := render.New(render.Options{
		TemplatesFS: rootFS,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/room-providers/{id}", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, re.Auto(c, RoomProvider{ID: 1, Name: "Ford"}))
	})

	w := httptest.New(app)
	res := w.HTML("/room-providers/1").Get()

	r.Contains(res.Body.String(), "SHOW: Ford")
}
