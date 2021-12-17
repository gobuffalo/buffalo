package buffalo

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/httptest"

	"github.com/stretchr/testify/require"
)

func Test_Session_SingleCookie(t *testing.T) {
	r := require.New(t)

	sessionName := "_test_session"
	a := New(Options{SessionName: sessionName})
	rr := render.New(render.Options{})

	a.GET("/", func(c Context) error {
		return c.Render(http.StatusCreated, rr.String(""))
	})

	w := httptest.New(a)
	res := w.HTML("/").Get()

	var sessionCookies []string
	for _, c := range res.Header().Values("Set-Cookie") {
		if strings.HasPrefix(c, sessionName) {
			sessionCookies = append(sessionCookies, c)
		}
	}

	r.Equal(1, len(sessionCookies))
}

func Test_Session_CustomValue(t *testing.T) {
	r := require.New(t)

	a := New(Options{})
	rr := render.New(render.Options{})

	// Root path sets a custom session value
	a.GET("/", func(c Context) error {
		c.Session().Set("example", "test")
		return c.Render(http.StatusCreated, rr.String(""))
	})
	// /session path prints custom session value as response
	a.GET("/session", func(c Context) error {
		sessionValue := c.Session().Get("example")
		return c.Render(http.StatusCreated, rr.String(fmt.Sprintf("%s", sessionValue)))
	})

	w := httptest.New(a)

	resSetSession := w.HTML("/").Get()

	// Create second request and set the cookie from the first response
	reqGetSession := w.HTML("/session")
	reqGetSession.Headers["Set-Cookie"] = resSetSession.Header().Values("Set-Cookie")[0]
	resGetSession := reqGetSession.Get()

	r.Equal(resGetSession.Body.String(), "test")
}
