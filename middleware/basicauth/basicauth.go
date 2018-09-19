package basicauth

import (
	basicauth "github.com/gobuffalo/mw-basicauth"
	"github.com/markbates/oncer"
)

var (
	// ErrNoCreds is returned when no basic auth credentials are defined
	ErrNoCreds = basicauth.ErrNoCreds

	// ErrAuthFail is returned when the client fails basic authentication
	ErrAuthFail = basicauth.ErrAuthFail
)

// Authorizer is used to authenticate the basic auth username/password.
// Should return true/false and/or an error.
//
// Deprecated: use github.com/gobuffalo/mw-basicauth#Authorizer instead.
type Authorizer = basicauth.Authorizer

// Middleware enables basic authentication
//
// Deprecated: use github.com/gobuffalo/mw-basicauth#Middleware instead.
var Middleware = basicauth.Middleware

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo/middleware/basicauth", "Use github.com/gobuffalo/mw-basicauth instead.")
}
