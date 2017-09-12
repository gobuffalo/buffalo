package middleware

import (
	"fmt"
	"runtime"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware/csrf"
)

var (
	// ErrNoReferer is returned when a HTTPS request provides an empty Referer
	// header.
	ErrNoReferer = csrf.ErrNoReferer
	// ErrBadReferer is returned when the scheme & host in the URL do not match
	// the supplied Referer header.
	ErrBadReferer = csrf.ErrBadReferer
	// ErrNoCSRFToken is returned if no CSRF token is supplied in the request.
	ErrNoCSRFToken = csrf.ErrNoToken
	// ErrBadCSRFToken is returned if the CSRF token in the request does not match
	// the token in the session, or is otherwise malformed.
	ErrBadCSRFToken = csrf.ErrBadToken
)

// CSRF is deprecated, and will be removed in v0.10.0. Use csrf.New instead.
var CSRF = func(next buffalo.Handler) buffalo.Handler {
	warningMsg := "middleware.CSRF is deprecated, and will be removed in v0.10.0. Use csrf.New instead."
	_, file, no, ok := runtime.Caller(1)
	if ok {
		warningMsg = fmt.Sprintf("%s Called from %s:%d", warningMsg, file, no)
	}
	return csrf.Middleware(next)
}
