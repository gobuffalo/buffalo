package buffalo

import (
	"testing"

	"github.com/markbates/willie"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_PanicHandler(t *testing.T) {
	app := New(Options{})
	app.GET("/string", func(c Context) error {
		panic("string boom!")
	})
	app.GET("/error", func(c Context) error {
		panic(errors.New("error boom!"))
	})

	table := []struct {
		path     string
		expected string
	}{
		{"/string", "string boom!"},
		{"/error", "error boom!"},
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
