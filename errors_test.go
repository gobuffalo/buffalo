package buffalo

import (
	"testing"

	"github.com/gobuffalo/httptest"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_defaultErrorHandler_SetsContentType(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(401, errors.New("boom"))
	})

	w := httptest.New(app)
	res := w.HTML("/").Get()
	r.Equal(401, res.Code)
	ct := res.Header().Get("content-type")
	r.Equal("text/html; charset=utf-8", ct)
}

func Test_defaultErrorHandler_JSON(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(401, errors.New("boom"))
	})

	w := httptest.New(app)
	res := w.JSON("/").Get()
	r.Equal(401, res.Code)
	ct := res.Header().Get("content-type")
	r.Equal("application/json", ct)
	b := res.Body.String()
	r.Contains(b, `"code":401`)
	r.Contains(b, `"error":"boom"`)
	r.Contains(b, `"trace":"`)
}

func Test_defaultErrorHandler_XML(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(401, errors.New("boom"))
	})

	w := httptest.New(app)
	res := w.XML("/").Get()
	r.Equal(401, res.Code)
	ct := res.Header().Get("content-type")
	r.Equal("application/xml", ct)
	b := res.Body.String()
	r.Contains(b, `<response code="401">`)
	r.Contains(b, `<error>boom</error>`)
	r.Contains(b, `<trace>`)
	r.Contains(b, `</trace>`)
	r.Contains(b, `</response>`)
}

func Test_PanicHandler(t *testing.T) {
	app := New(Options{})
	app.GET("/string", func(c Context) error {
		panic("string boom")
	})
	app.GET("/error", func(c Context) error {
		panic(errors.New("error boom"))
	})

	table := []struct {
		path     string
		expected string
	}{
		{"/string", "string boom"},
		{"/error", "error boom"},
	}

	const stack = `github.com/gobuffalo/buffalo.(*App).PanicHandler`

	w := httptest.New(app)
	for _, tt := range table {
		t.Run(tt.path, func(st *testing.T) {
			r := require.New(st)

			res := w.HTML(tt.path).Get()
			r.Equal(500, res.Code)

			body := res.Body.String()
			r.Contains(body, tt.expected)
			r.Contains(body, stack)
		})
	}
}

func Test_defaultErrorMiddleware(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	var x string
	var ok bool
	app.ErrorHandlers[422] = func(code int, err error, c Context) error {
		x, ok = c.Value("T").(string)
		c.Response().WriteHeader(code)
		c.Response().Write([]byte(err.Error()))
		return nil
	}
	app.Use(func(next Handler) Handler {
		return func(c Context) error {
			c.Set("T", "t")
			return c.Error(422, errors.New("boom"))
		}
	})
	app.GET("/", func(c Context) error {
		return nil
	})

	w := httptest.New(app)
	res := w.HTML("/").Get()
	r.Equal(422, res.Code)
	r.True(ok)
	r.Equal("t", x)
}
