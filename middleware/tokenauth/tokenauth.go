package tokenauth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

var (
	// ErrTokenInvalid is returned when the token provided is invalid
	ErrTokenInvalid = errors.New("token invalid")
	// ErrNoToken is returned if no token is supplied in the request.
	ErrNoToken = errors.New("token not found in request")
	// ErrBadToken is returned if the token sign method in the request
	// does not match the signing method used
	ErrBadSigningMethod = errors.New("unexpected signing method")
)

func Middleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		authString := c.Request().Header.Get("Authorization")
		SecretKey := envy.Get("JWT_SECRET", "secret")

		tokenString, err := getJwtToken(authString)
		if err != nil {
			return c.Error(http.StatusUnauthorized, err)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrBadSigningMethod
			}
			return []byte(SecretKey), nil
		})
		if err != nil {
			return c.Error(http.StatusUnauthorized, err)
		}
		c.Set("claims", token.Claims)
		// do some work before calling the next handler
		err = next(c)
		// do some work after calling the next handler
		return err
	}
}

func getJwtToken(authString string) (string, error) {
	if authString == "" {
		return "", ErrNoToken
	}
	splitToken := strings.Split(authString, "Bearer ")
	if len(splitToken) != 2 {
		return "", ErrTokenInvalid
	}
	tokenString := splitToken[1]
	return tokenString, nil
}
