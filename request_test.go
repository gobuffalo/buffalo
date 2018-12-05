package buffalo

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gobuffalo/httptest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func Test_Request_MultipleReads(t *testing.T) {
	r := require.New(t)
	var reads []string

	h := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, err := ioutil.ReadAll(req.Body)
		r.NoError(err)
		reads = append(reads, string(b))
	})
	app := New(Options{
		PreHandlers: []http.Handler{h},
	})

	app.Use(func(next Handler) Handler {
		return func(c Context) error {
			b, err := ioutil.ReadAll(c.Request().Body)
			if err != nil {
				return errors.WithStack(err)
			}
			reads = append(reads, string(b))
			return next(c)
		}
	})
	app.POST("/", func(c Context) error {
		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return errors.WithStack(err)
		}
		reads = append(reads, string(b))
		return nil
	})

	w := httptest.New(app)
	w.JSON("/").Post(map[string]string{"foo": "foo"})
	r.Len(reads, 3)

	foo := `{"foo":"foo"}`
	r.Equal([]string{foo, foo, foo}, reads)
}
