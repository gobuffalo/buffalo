package tokenauth_test

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
	"net/http"
	"testing"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware/tokenauth"
	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func app() *buffalo.App {
	h := func(c buffalo.Context) error {
		return c.Render(200, nil)
	}

	a := buffalo.New(buffalo.Options{})
	a.Use(tokenauth.Middleware)
	a.GET("/", h)
	return a
}

func TestTokenMiddleware(t *testing.T) {
	r := require.New(t)
	w := willie.New(app())

	// Missing Authorization
	res := w.Request("/").Get()
	r.Equal(http.StatusUnauthorized, res.Code)

	// invalid token
	req := w.Request("/")
	req.Headers["Authorization"] = "badcreds"
	res = req.Get()
	r.Equal(http.StatusUnauthorized, res.Code)
	r.Contains(res.Body.String(), "token invalid")

	// expired token
	SecretKey := envy.Get("JWT_SECRET", "secret")
	claims := jwt.MapClaims{}
	claims["sub"] = "1234567890"
	claims["exp"] = time.Now().Add(-time.Minute * 5).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(SecretKey))
	req.Headers["Authorization"] = fmt.Sprintf("Bearer %s", tokenString)
	res = req.Get()
	fmt.Println(res.Body.String())
	r.Equal(http.StatusUnauthorized, res.Code)
	r.Contains(res.Body.String(), "Token is expired")

	// valid token
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ = token.SignedString([]byte(SecretKey))
	req.Headers["Authorization"] = fmt.Sprintf("Bearer %s", tokenString)
	res = req.Get()
	r.Equal(http.StatusOK, res.Code)
}
