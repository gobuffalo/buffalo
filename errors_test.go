package buffalo

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/logger"
	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"
)

//testLoggerHook is useful to test whats being logged.
type testLoggerHook struct {
	errors []*logrus.Entry
}

func (lh *testLoggerHook) Fire(entry *logrus.Entry) error {
	lh.errors = append(lh.errors, entry)
	return nil
}

func (lh *testLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
	}
}

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

	testHook := &testLoggerHook{}
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.AddHook(testHook)
	log := logger.Logrus{
		FieldLogger: l,
	}
	app.Logger = log

	w := httptest.New(app)
	res := w.HTML("/").Get()
	r.Equal(http.StatusUnauthorized, res.Code)
	r.Equal(http.StatusUnauthorized, testHook.errors[0].Data["status"])
}

func Test_defaultErrorHandler_JSON(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(http.StatusUnauthorized, fmt.Errorf("boom"))
	})

	w := httptest.New(app)
	res := w.JSON("/").Get()
	r.Equal(http.StatusUnauthorized, res.Code)
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
		return c.Error(http.StatusUnauthorized, fmt.Errorf("boom"))
	})

	w := httptest.New(app)
	res := w.XML("/").Get()
	r.Equal(http.StatusUnauthorized, res.Code)
	ct := res.Header().Get("content-type")
	r.Equal("text/xml", ct)
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

			res := w.HTML(tt.path).Get()
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
