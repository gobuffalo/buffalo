package buffalo

import (
	"net/url"
	"testing"

	"github.com/markbates/buffalo/render"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_MethodOverride(t *testing.T) {
	r := require.New(t)

	a := Automatic(Options{})
	a.PUT("/", func(c Context) error {
		return c.Render(200, render.String("you put me!"))
	})

	w := willie.New(a)
	res := w.Request("/").Post(url.Values{"_method": []string{"PUT"}})
	r.Equal(200, res.Code)
	r.Equal("you put me!", res.Body.String())
}
