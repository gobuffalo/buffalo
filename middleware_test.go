package buffalo

import (
	"net/http"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"
	"github.com/stretchr/testify/require"
)

// Test_App_Use tests that middleware gets added
func Test_App_Use(t *testing.T) {
	r := require.New(t)

	log := []string{}
	a := New(Options{})
	a.Use(func(h Handler) Handler {
		return func(c Context) error {
			log = append(log, "start")
			err := h(c)
			log = append(log, "end")
			return err
		}
	})

	a.GET("/", func(c Context) error {
		log = append(log, "handler")
		return nil
	})

	w := httptest.New(a)
	w.HTML("/").Get()
	r.Len(log, 3)
	r.Equal([]string{"start", "handler", "end"}, log)
}

// Test_Middleware_Replace tests that middleware gets added
func Test_Middleware_Replace(t *testing.T) {
	r := require.New(t)

	log := []string{}
	a := New(Options{})
	mw1 := func(h Handler) Handler {
		return func(c Context) error {
			log = append(log, "m1 start")
			err := h(c)
			log = append(log, "m1 end")
			return err
		}
	}
	mw2 := func(h Handler) Handler {
		return func(c Context) error {
			log = append(log, "m2 start")
			err := h(c)
			log = append(log, "m2 end")
			return err
		}
	}
	a.Use(mw1)
	a.Middleware.Replace(mw1, mw2)

	a.GET("/", func(c Context) error {
		log = append(log, "handler")
		return nil
	})

	w := httptest.New(a)
	w.HTML("/").Get()
	r.Len(log, 3)
	r.Equal([]string{"m2 start", "handler", "m2 end"}, log)
}

// Test_Middleware_Skip tests that middleware gets skipped
func Test_Middleware_Skip(t *testing.T) {
	r := require.New(t)

	log := []string{}
	a := New(Options{})
	mw1 := func(h Handler) Handler {
		return func(c Context) error {
			log = append(log, "mw1 start")
			err := h(c)
			log = append(log, "mw1 end")
			return err
		}
	}
	mw2 := func(h Handler) Handler {
		return func(c Context) error {
			log = append(log, "mw2 start")
			err := h(c)
			log = append(log, "mw2 end")
			return err
		}
	}
	a.Use(mw1)
	a.Use(mw2)

	h1 := func(c Context) error {
		log = append(log, "h1")
		return nil
	}
	h2 := func(c Context) error {
		log = append(log, "h2")
		return nil
	}

	a.GET("/h1", h1)
	a.GET("/h2", h2)

	a.Middleware.Skip(mw2, h2)

	w := httptest.New(a)

	w.HTML("/h2").Get()
	r.Len(log, 3)
	r.Equal([]string{"mw1 start", "h2", "mw1 end"}, log)

	log = []string{}
	w.HTML("/h1").Get()
	r.Len(log, 5)
	r.Equal([]string{"mw1 start", "mw2 start", "h1", "mw2 end", "mw1 end"}, log)
}

type carsResource struct {
	Resource
}

func (ur *carsResource) Show(c Context) error {
	return c.Render(http.StatusOK, render.String("show"))
}

func (ur *carsResource) List(c Context) error {
	return c.Render(http.StatusOK, render.String("list"))
}

// Test_Middleware_Skip tests that middleware gets skipped
func Test_Middleware_Skip_Resource(t *testing.T) {
	r := require.New(t)

	log := []string{}
	mw1 := func(h Handler) Handler {
		return func(c Context) error {
			log = append(log, "mw1 start")
			err := h(c)
			log = append(log, "mw1 end")
			return err
		}
	}

	a := New(Options{})
	var cr Resource = &carsResource{}
	g := a.Resource("/autos", cr)
	g.Use(mw1)

	var ur Resource = &carsResource{}
	g = a.Resource("/cars", ur)
	g.Use(mw1)

	// fmt.Println("set up skip")
	g.Middleware.Skip(mw1, ur.Show)

	w := httptest.New(a)

	// fmt.Println("make autos call")
	log = []string{}
	res := w.HTML("/autos/1").Get()
	r.Len(log, 2)
	r.Equal("show", res.Body.String())

	// fmt.Println("make list call")
	log = []string{}
	res = w.HTML("/cars").Get()
	r.Len(log, 2)
	r.Equal([]string{"mw1 start", "mw1 end"}, log)
	r.Equal("list", res.Body.String())

	// fmt.Println("make show call")
	log = []string{}
	res = w.HTML("/cars/1").Get()
	r.Len(log, 0)
	r.Equal("show", res.Body.String())

}

// Test_Middleware_Clear confirms that middle gets cleared
func Test_Middleware_Clear(t *testing.T) {
	r := require.New(t)
	mws := newMiddlewareStack()
	mw := func(h Handler) Handler { return h }
	mws.Use(mw)
	mws.Skip(mw, voidHandler)

	r.Len(mws.stack, 1)
	r.Len(mws.skips, 1)

	mws.Clear()

	r.Len(mws.stack, 0)
	r.Len(mws.skips, 0)
}

func Test_Middleware_Remove(t *testing.T) {
	r := require.New(t)
	log := []string{}

	mw1 := func(h Handler) Handler {
		log = append(log, "mw1")
		return h
	}

	mw2 := func(h Handler) Handler {
		log = append(log, "mw2")
		return h
	}

	a := New(Options{})
	a.Use(mw2)
	a.Use(mw1)

	var cr Resource = &carsResource{}
	g := a.Resource("/autos", cr)
	g.Middleware.Remove(mw2)

	a.Resource("/all_log_autos", cr)
	w := httptest.New(a)

	ng := a.Resource("/no_log_autos", cr)
	ng.Middleware.Remove(mw1, mw2)

	_ = w.HTML("/autos/1").Get()
	r.Len(log, 1)
	r.Equal("mw1", log[0])

	log = []string{}
	_ = w.HTML("/all_log_autos/1").Get()
	r.Len(log, 2)
	r.Contains(log, "mw2")
	r.Contains(log, "mw1")

	log = []string{}
	_ = w.HTML("/no_log_autos/1").Get()
	r.Len(log, 0)
}
