package buffalo

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gobuffalo/httptest"

	"github.com/stretchr/testify/require"
)

func Test_defaultErrorHandler_SetsContentType(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("boom"))
	})

	w := httptest.New(app)
	res := w.HTML("/").Get()
	r.Equal(http.StatusUnauthorized, res.Code)
	ct := res.Header().Get("content-type")
	r.Equal("text/html; charset=utf-8", ct)
}

func Test_defaultErrorHandler_Logger(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("boom"))
	})

	w := httptest.New(app)
	res := w.HTML("/").Get()
	r.Equal(http.StatusUnauthorized, res.Code)
}

func Test_defaultErrorHandler_JSON_development(t *testing.T) {
	testDefaultErrorHandler(t, "application/json", "development")
}

func Test_defaultErrorHandler_XML_development(t *testing.T) {
	testDefaultErrorHandler(t, "text/xml", "development")
}

func Test_defaultErrorHandler_JSON_staging(t *testing.T) {
	testDefaultErrorHandler(t, "application/json", "staging")
}

func Test_defaultErrorHandler_XML_staging(t *testing.T) {
	testDefaultErrorHandler(t, "text/xml", "staging")
}

func Test_defaultErrorHandler_JSON_production(t *testing.T) {
	testDefaultErrorHandler(t, "application/json", "production")
}

func Test_defaultErrorHandler_XML_production(t *testing.T) {
	testDefaultErrorHandler(t, "text/xml", "production")
}

func testDefaultErrorHandler(t *testing.T, contentType, env string) {
	r := require.New(t)
	app := New(Options{})
	app.Env = env
	app.GET("/", func(c Context) error {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("boom"))
	})

	w := httptest.New(app)
	var res *httptest.Response
	if contentType == "application/json" {
		res = w.JSON("/").Get().Response
	} else {
		res = w.XML("/").Get().Response
	}
	r.Equal(http.StatusUnauthorized, res.Code)
	ct := res.Header().Get("content-type")
	r.Equal(contentType, ct)
	b := res.Body.String()

	if env == "development" {
		if contentType == "text/xml" {
			r.Contains(b, `<response code="401">`)
			r.Contains(b, `<error>boom</error>`)
			r.Contains(b, `<trace>`)
			r.Contains(b, `</trace>`)
			r.Contains(b, `</response>`)
		} else {
			r.Contains(b, `"code":401`)
			r.Contains(b, `"error":"boom"`)
			r.Contains(b, `"trace":"`)
		}
	} else {
		if contentType == "text/xml" {
			r.Contains(b, `<response code="401">`)
			r.Contains(b, fmt.Sprintf(`<error>%s</error>`, http.StatusText(http.StatusUnauthorized)))
			r.NotContains(b, `<trace>`)
			r.NotContains(b, `</trace>`)
			r.Contains(b, `</response>`)
		} else {
			r.Contains(b, `"code":401`)
			r.Contains(b, fmt.Sprintf(`"error":"%s"`, http.StatusText(http.StatusUnauthorized)))
			r.NotContains(b, `"trace":"`)
		}
	}
}

func Test_defaultErrorHandler_nil_error(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(http.StatusInternalServerError, nil)
	})

	w := httptest.New(app)
	res := w.JSON("/").Get()
	r.Equal(http.StatusInternalServerError, res.Code)
}

func Test_PanicHandler(t *testing.T) {
	app := New(Options{})
	app.GET("/string", func(c Context) error {
		panic("string boom")
	})
	app.GET("/error", func(c Context) error {
		panic(fmt.Errorf("error boom"))
	})

	table := []struct {
		path     string
		expected string
	}{
		{"/string", "string boom"},
		{"/error", "error boom"},
	}

	const stack = `github.com/gobuffalo/buffalo.Test_PanicHandler`

	w := httptest.New(app)
	for _, tt := range table {
		t.Run(tt.path, func(st *testing.T) {
			r := require.New(st)

			res := w.HTML("%s", tt.path).Get()
			r.Equal(http.StatusInternalServerError, res.Code)

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
	app.ErrorHandlers[http.StatusUnprocessableEntity] = func(code int, err error, c Context) error {
		x, ok = c.Value("T").(string)
		c.Response().WriteHeader(code)
		c.Response().Write([]byte(err.Error()))
		return nil
	}
	app.Use(func(next Handler) Handler {
		return func(c Context) error {
			c.Set("T", "t")
			return c.Error(http.StatusUnprocessableEntity, fmt.Errorf("boom"))
		}
	})
	app.GET("/", func(c Context) error {
		return nil
	})

	w := httptest.New(app)
	res := w.HTML("/").Get()
	r.Equal(http.StatusUnprocessableEntity, res.Code)
	r.True(ok)
	r.Equal("t", x)
}

func Test_SetErrorMiddleware(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.ErrorHandlers.Default(func(code int, err error, c Context) error {
		res := c.Response()
		res.WriteHeader(http.StatusTeapot)
		res.Write([]byte("i'm a teapot"))
		return nil
	})
	app.GET("/", func(c Context) error {
		return c.Error(http.StatusUnprocessableEntity, fmt.Errorf("boom"))
	})

	w := httptest.New(app)
	res := w.HTML("/").Get()
	r.Equal(http.StatusTeapot, res.Code)
	r.Equal("i'm a teapot", res.Body.String())
}
