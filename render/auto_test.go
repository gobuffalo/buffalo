package render_test

import (
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/packd"
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
		return c.Render(200, re.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
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
		return c.Render(200, re.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
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
		return c.Render(200, re.Auto(c, []string{"Honda", "Toyota", "Ford", "Chevy"}))
	})

	w := httptest.New(app)

	res := w.XML("/cars").Get()
	r.Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<string>Honda</string>\n<string>Toyota</string>\n<string>Ford</string>\n<string>Chevy</string>", strings.TrimSpace(res.Body.String()))
}

func Test_Auto_HTML_List(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	err := box.AddString("cars/index.html", "INDEX: <%= len(cars) %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, Cars{
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

	box := packd.NewMemoryBox()
	err := box.AddString("people/index.html", "INDEX: <%= len(people) %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/people", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, People{
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

	box := packd.NewMemoryBox()
	err := box.AddString("cars/show.html", "Show: <%= car.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars/{id}", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Get()
	r.Contains(res.Body.String(), "Show: Honda")
	r.NoError(err)
}

func Test_Auto_HTML_New(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	err := box.AddString("cars/new.html", "New: <%= car.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars/new", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/new").Get()
	r.Contains(res.Body.String(), "New: Honda")
}

func Test_Auto_HTML_Create(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	err := box.AddString("cars/new.html", "New: <%= car.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.POST("/cars", func(c buffalo.Context) error {
		return c.Render(201, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars").Post(nil)
	r.Contains(res.Body.String(), "New: Honda")
}

func Test_Auto_HTML_Create_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	app.POST("/cars", func(c buffalo.Context) error {
		return c.Render(201, render.Auto(c, Car{
			ID:   1,
			Name: "Honda",
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars").Post(nil)
	r.Equal("/cars/1", res.Location())
	r.Equal(302, res.Code)
}

func Test_Auto_HTML_Create_Redirect_Error(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	err := box.AddString("cars/new.html", "Create: <%= car.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.POST("/cars", func(c buffalo.Context) error {
		b := Car{
			Name: "Honda",
		}
		return c.Render(422, re.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars").Post(nil)
	r.Equal(422, res.Code)
	r.Contains(res.Body.String(), "Create: Honda")
}

func Test_Auto_HTML_Create_Nested_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	admin := app.Group("/admin")
	admin.POST("/cars", func(c buffalo.Context) error {
		return c.Render(201, render.Auto(c, Car{
			ID:   1,
			Name: "Honda",
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/admin/cars").Post(nil)
	r.Equal("/admin/cars/1", res.Location())
	r.Equal(302, res.Code)
}

func Test_Auto_HTML_Destroy_Nested_Redirect(t *testing.T) {
	r := require.New(t)

	app := buffalo.New(buffalo.Options{})
	admin := app.Group("/admin")
	admin.DELETE("/cars", func(c buffalo.Context) error {
		return c.Render(200, render.Auto(c, Car{
			ID:   1,
			Name: "Honda",
		}))
	})

	w := httptest.New(app)
	res := w.HTML("/admin/cars").Delete()
	r.Equal("/admin/cars", res.Location())
	r.Equal(302, res.Code)
}

func Test_Auto_HTML_Edit(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	err := box.AddString("cars/edit.html", "Edit: <%= car.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/cars/{id}/edit", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, Car{Name: "Honda"}))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1/edit").Get()
	r.Contains(res.Body.String(), "Edit: Honda")
}

func Test_Auto_HTML_Update(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	err := box.AddString("cars/edit.html", "Update: <%= car.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.PUT("/cars/{id}", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, Car{Name: "Honda"}))
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
		return c.Render(200, render.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Put(nil)
	r.Equal("/cars/1", res.Location())
	r.Equal(302, res.Code)
}

func Test_Auto_HTML_Update_Redirect_Error(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	err := box.AddString("cars/edit.html", "Update: <%= car.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.PUT("/cars/{id}", func(c buffalo.Context) error {
		b := Car{
			ID:   1,
			Name: "Honda",
		}
		return c.Render(422, re.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Put(nil)
	r.Equal(422, res.Code)
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
		return c.Render(200, render.Auto(c, b))
	})

	w := httptest.New(app)
	res := w.HTML("/cars/1").Delete()
	r.Equal("/cars/", res.Location())
	r.Equal(302, res.Code)
}

func Test_Auto_HTML_List_Plural_MultiWord(t *testing.T) {
	r := require.New(t)

	type RoomProvider struct {
		Name string
	}

	type RoomProviders []RoomProvider

	box := packd.NewMemoryBox()
	err := box.AddString("room_providers/index.html", "INDEX: <%= len(roomProviders) %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/room_providers", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, RoomProviders{
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

	box := packd.NewMemoryBox()
	err := box.AddString("room_providers/index.html", "INDEX: <%= len(roomProviders) %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/room-providers", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, RoomProviders{
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

	box := packd.NewMemoryBox()
	err := box.AddString("room_providers/show.html", "SHOW: <%= roomProvider.Name %>")
	r.NoError(err)

	re := render.New(render.Options{
		TemplatesBox: box,
	})

	app := buffalo.New(buffalo.Options{})
	app.GET("/room-providers/{id}", func(c buffalo.Context) error {
		return c.Render(200, re.Auto(c, RoomProvider{ID: 1, Name: "Ford"}))
	})

	w := httptest.New(app)
	res := w.HTML("/room-providers/1").Get()

	r.Contains(res.Body.String(), "SHOW: Ford")
}
