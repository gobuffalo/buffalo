package buffalo

import (
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/markbates/willie"
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

	w := willie.New(a)
	w.Request("/").Get()
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

	w := willie.New(a)
	w.Request("/").Get()
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

	w := willie.New(a)

	w.Request("/h2").Get()
	r.Len(log, 3)
	r.Equal([]string{"mw1 start", "h2", "mw1 end"}, log)

	log = []string{}
	w.Request("/h1").Get()
	r.Len(log, 5)
	r.Equal([]string{"mw1 start", "mw2 start", "h1", "mw2 end", "mw1 end"}, log)
}

type carsResource struct {
	Resource
}

func (ur *carsResource) Show(c Context) error {
	return c.Render(200, render.String("show"))
}

func (ur *carsResource) List(c Context) error {
	return c.Render(200, render.String("list"))
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

	w := willie.New(a)

	// fmt.Println("make autos call")
	log = []string{}
	res := w.Request("/autos/1").Get()
	r.Len(log, 2)
	r.Equal("show", res.Body.String())

	// fmt.Println("make list call")
	log = []string{}
	res = w.Request("/cars").Get()
	r.Len(log, 2)
	r.Equal([]string{"mw1 start", "mw1 end"}, log)
	r.Equal("list", res.Body.String())

	// fmt.Println("make show call")
	log = []string{}
	res = w.Request("/cars/1").Get()
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
