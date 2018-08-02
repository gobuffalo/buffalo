package buffalo

import (
	"testing"

	"github.com/markbates/willie"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_defaultErrorHandler_SetsContentType(t *testing.T) {
	r := require.New(t)
	app := New(Options{})
	app.GET("/", func(c Context) error {
		return c.Error(401, errors.New("boom"))
	})

	w := willie.New(app)
	res := w.HTML("/").Get()
	r.Equal(401, res.Code)
	ct := res.Header().Get("content-type")
	r.Equal("text/html", ct)
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

	w := willie.New(app)
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
