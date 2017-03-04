package middleware_test

import (
	"testing"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/render"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

type csrfForm struct {
	AuthenticityToken string `form:"authenticity_token"`
}

func ctCSRFApp() *buffalo.App {
	h := func(c buffalo.Context) error {
		return c.Render(200, render.String(c.(*buffalo.DefaultContext).Data()["authenticity_token"].(string)))
	}
	a := buffalo.Automatic(buffalo.Options{})
	a.GET("/csrf", middleware.EnableCSRF()(h))
	a.POST("/csrf", middleware.EnableCSRF()(h))
	return a
}

func Test_CSRFOnIdempotentAction(t *testing.T) {
	r := require.New(t)

	w := willie.New(ctCSRFApp())
	res := w.Request("/csrf").Get()
	r.Equal(200, res.Code)
}

func Test_CSRFOnEditingAction(t *testing.T) {
	r := require.New(t)

	w := willie.New(ctCSRFApp())

	// Test missing token case
	res := w.Request("/csrf").Post("")
	r.Equal(500, res.Code)
	r.Contains(res.Body.String(), "CSRF token not found in request")

	// Test provided bad token through Header case
	req := w.Request("/csrf")
	req.Headers["X-CSRF-Token"] = "test-token"
	res = req.Post("")
	r.Equal(500, res.Code)
	r.Contains(res.Body.String(), "CSRF token not found in request")

	// Test provided good token through Header case
	res = w.Request("/csrf").Get()
	r.Equal(200, res.Code)
	token := res.Body.String()

	req = w.Request("/csrf")
	req.Headers["X-CSRF-Token"] = token
	res = req.Post("")
	r.Equal(200, res.Code)

	// Test provided good token through form case
	res = w.Request("/csrf").Get()
	r.Equal(200, res.Code)
	token = res.Body.String()

	req = w.Request("/csrf")
	res = req.Post(csrfForm{AuthenticityToken: token})
	r.Equal(200, res.Code)
}
