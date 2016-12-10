package buffalo

import (
	"testing"

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
