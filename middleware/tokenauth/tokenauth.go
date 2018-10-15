// Package tokenauth provides jwt token authorisation middleware
// supports HMAC, RSA, ECDSA, RSAPSS algorithms
// uses github.com/dgrijalva/jwt-go for jwt implementation
//
// Setting Up tokenauth middleware
//
// Using tokenauth with defaults
//  app.Use(tokenauth.New(tokenauth.Options{}))
// Specifying Signing method for JWT
//  app.Use(tokenauth.New(tokenauth.Options{
//      SignMethod: jwt.SigningMethodRS256,
//  }))
// By default the Key used is loaded from the JWT_SECRET or JWT_PUBLIC_KEY env variable depending
// on the SigningMethod used. However you can retrive the key from a different source.
//  app.Use(tokenauth.New(tokenauth.Options{
//      GetKey: func(jwt.SigningMethod) (interface{}, error) {
//           // Your Implementation here ...
//      },
//  }))
//
//
// Creating a new token
//
// This can be referred from the underlying JWT package being used https://github.com/dgrijalva/jwt-go
//
// Example
//  claims := jwt.MapClaims{}
//  claims["userid"] = "123"
//  claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
//  // add more claims
//  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//  tokenString, err := token.SignedString([]byte(SecretKey))
//
//
// Getting Claims from JWT token from buffalo context
//
// Example of retriving username from claims (this step is same regardless of the signing method used)
//  claims := c.Value("claims").(jwt.MapClaims)
//  username := claims["username"].(string)
package tokenauth

import (
	tokenauth "github.com/gobuffalo/mw-tokenauth"
	"github.com/markbates/oncer"
)

var (
	// ErrTokenInvalid is returned when the token provided is invalid
	ErrTokenInvalid = tokenauth.ErrTokenInvalid
	// ErrNoToken is returned if no token is supplied in the request.
	ErrNoToken = tokenauth.ErrNoToken
	// ErrBadSigningMethod is returned if the token sign method in the request
	// does not match the signing method used
	ErrBadSigningMethod = tokenauth.ErrBadSigningMethod
)

// Options for the JWT middleware
//
// Deprecated: use github.com/gobuffalo/mw-tokenauth#Options instead.
type Options = tokenauth.Options

// New enables jwt token verification if no Sign method is provided,
// by default uses HMAC
//
// Deprecated: use github.com/gobuffalo/mw-tokenauth#New instead.
var New = tokenauth.New

// GetHMACKey gets secret key from env
//
// Deprecated: use github.com/gobuffalo/mw-tokenauth#GetHMACKey instead.
var GetHMACKey = tokenauth.GetHMACKey

// GetKeyRSA gets the public key file location from env and returns rsa.PublicKey
//
// Deprecated: use github.com/gobuffalo/mw-tokenauth#GetKeyRSA instead.
var GetKeyRSA = tokenauth.GetKeyRSA

// GetKeyRSAPSS uses GetKeyRSA() since both requires rsa.PublicKey
//
// Deprecated: use github.com/gobuffalo/mw-tokenauth#GetKeyRSAPSS instead.
var GetKeyRSAPSS = tokenauth.GetKeyRSAPSS

// GetKeyECDSA gets the public.pem file location from env and returns ecdsa.PublicKey
//
// Deprecated: use github.com/gobuffalo/mw-tokenauth#GetKeyECDSA instead.
var GetKeyECDSA = tokenauth.GetKeyECDSA

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo/middleware/tokenauth", "Use github.com/gobuffalo/mw-tokenauth instead.")
}
