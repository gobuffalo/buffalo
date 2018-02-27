// Package tokenauth provides jwt token authorisation middleware
// supports HMAC, RSA, ECDSA, RSAPSS algorithms
// uses github.com/dgrijalva/jwt-go for jwt implementation
package tokenauth

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

var (
	// ErrTokenInvalid is returned when the token provided is invalid
	ErrTokenInvalid = errors.New("token invalid")
	// ErrNoToken is returned if no token is supplied in the request.
	ErrNoToken = errors.New("token not found in request")
	// ErrBadSigningMethod is returned if the token sign method in the request
	// does not match the signing method used
	ErrBadSigningMethod = errors.New("unexpected signing method")
)

// Options for the JWT middleware
type Options struct {
	SignMethod jwt.SigningMethod
	GetKey     func(jwt.SigningMethod) (interface{}, error)
}

// New enables jwt token verification if no Sign method is provided,
// by default uses HMAC
func New(options Options) buffalo.MiddlewareFunc {
	// set sign method to HMAC if not provided
	if options.SignMethod == nil {
		options.SignMethod = jwt.SigningMethodHS256
	}
	if options.GetKey == nil {
		options.GetKey = selectGetKeyFunc(options.SignMethod)
	}
	// get key for validation
	key, err := options.GetKey(options.SignMethod)
	// if error on getting key exit.
	if err != nil {
		log.Fatal(errors.Wrap(err, "couldn't get key"))
	}
	return func(next buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			// get Authorisation header value
			authString := c.Request().Header.Get("Authorization")

			tokenString, err := getJwtToken(authString)
			// if error on getting the token, return with status unauthorized
			if err != nil {
				return c.Error(http.StatusUnauthorized, err)
			}

			// validating and parsing the tokenString
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validating if algorithm used for signing is same as the algorithm in token
				if token.Method.Alg() != options.SignMethod.Alg() {
					return nil, ErrBadSigningMethod
				}
				return key, nil
			})
			// if error validating jwt token, return with status unauthorized
			if err != nil {
				return c.Error(http.StatusUnauthorized, err)
			}

			// set the claims as context parameter.
			// so that the actions can use the claims from jwt token
			c.Set("claims", token.Claims)
			// calling next handler
			err = next(c)

			return err
		}
	}
}

// selectGetKeyFunc is an helper function to choose the GetKey function
// according to the Signing method used
func selectGetKeyFunc(method jwt.SigningMethod) func(jwt.SigningMethod) (interface{}, error) {
	switch method.(type) {
	case *jwt.SigningMethodRSA:
		return GetKeyRSA
	case *jwt.SigningMethodECDSA:
		return GetKeyECDSA
	case *jwt.SigningMethodRSAPSS:
		return GetKeyRSAPSS
	default:
		return GetHMACKey
	}
}

// GetHMACKey gets secret key from env
func GetHMACKey(jwt.SigningMethod) (interface{}, error) {
	key, err := envy.MustGet("JWT_SECRET")
	return []byte(key), err
}

// GetKeyRSA gets the public key file location from env and returns rsa.PublicKey
func GetKeyRSA(jwt.SigningMethod) (interface{}, error) {
	key, err := envy.MustGet("JWT_PUBLIC_KEY")
	if err != nil {
		return nil, err
	}
	keyData, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(keyData)
}

// GetKeyRSAPSS uses GetKeyRSA() since both requires rsa.PublicKey
func GetKeyRSAPSS(signingMethod jwt.SigningMethod) (interface{}, error) {
	return GetKeyRSA(signingMethod)
}

// GetKeyECDSA gets the public.pem file location from env and returns ecdsa.PublicKey
func GetKeyECDSA(jwt.SigningMethod) (interface{}, error) {
	key, err := envy.MustGet("JWT_PUBLIC_KEY")
	if err != nil {
		return nil, err
	}
	keyData, err := ioutil.ReadFile(key)
	if err != nil {
		return nil, err
	}
	return jwt.ParseECPublicKeyFromPEM(keyData)
}

// getJwtToken gets the token from the Authorisation header
// removes the Bearer part from the authorisation header value.
// returns No token error if Token is not found
// returns Token Invalid error if the token value cannot be obtained by removing `Bearer `
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
