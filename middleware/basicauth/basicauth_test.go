package basicauth_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware/basicauth"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func app() *buffalo.App {
	h := func(c buffalo.Context) error {
		return c.Render(200, nil)
	}
	auth := func(c buffalo.Context, u, p string) (bool, error) {
		return (u == "tester" && p == "pass123"), nil
	}
	a := buffalo.New(buffalo.Options{})
	a.Use(basicauth.Middleware(auth))
	a.GET("/", h)
	return a
}

func TestBasicAuth(t *testing.T) {
	r := require.New(t)

	w := willie.New(app())

	authfail := "invalid basic auth"

	// missing authorization
	res := w.Request("/").Get()
	r.Equal(401, res.Code)
	r.Contains(res.Header().Get("WWW-Authenticate"), `Basic realm="Basic Authentication"`)
	r.Contains(res.Body.String(), "Unauthorized")

	// bad header value, not Basic
	req := w.Request("/")
	req.Headers["Authorization"] = "badcreds"
	res = req.Get()
	r.Equal(401, res.Code)
	r.Contains(res.Body.String(), "Unauthorized")

	// bad cred values
	req = w.Request("/")
	req.Headers["Authorization"] = "bad creds"
	res = req.Get()
	r.Equal(500, res.Code)
	r.Contains(res.Body.String(), authfail)

	creds := base64.StdEncoding.EncodeToString([]byte("badcredvalue"))

	// invalid cred values in authorization
	req = w.Request("/")
	req.Headers["Authorization"] = fmt.Sprintf("Basic %s", creds)
	res = req.Get()
	r.Equal(500, res.Code)
	r.Contains(res.Body.String(), authfail)

	creds = base64.StdEncoding.EncodeToString([]byte("tester:pass123"))

	// valid cred values
	req = w.Request("/")
	req.Headers["Authorization"] = fmt.Sprintf("Basic %s", creds)
	res = req.Get()
	r.Equal(200, res.Code)
}
