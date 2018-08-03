package csrf

import (
	"fmt"

	"github.com/gobuffalo/mw-csrf"
)

const (
	// CSRF token length in bytes.
	tokenLength int    = 32
	tokenKey    string = "authenticity_token"
)

var (
	// The name value used in form fields.
	fieldName = tokenKey

	// The HTTP request header to inspect
	headerName = "X-CSRF-Token"

	// Idempotent (safe) methods as defined by RFC7231 section 4.2.2.
	safeMethods = []string{"GET", "HEAD", "OPTIONS", "TRACE"}
	htmlTypes   = []string{"html", "form", "plain", "*/*"}
)

var (
	// ErrNoReferer is returned when a HTTPS request provides an empty Referer
	// header.
	ErrNoReferer = csrf.ErrNoReferer
	// ErrBadReferer is returned when the scheme & host in the URL do not match
	// the supplied Referer header.
	ErrBadReferer = csrf.ErrBadReferer
	// ErrNoToken is returned if no CSRF token is supplied in the request.
	ErrNoToken = csrf.ErrNoToken
	// ErrBadToken is returned if the CSRF token in the request does not match
	// the token in the session, or is otherwise malformed.
	ErrBadToken = csrf.ErrBadToken
)

// New enable CSRF protection on routes using this middleware.
// This middleware is adapted from gorilla/csrf
//
// Deprecated: use github.com/gobuffalo/mw-csrf#New instead.
var New = csrf.New

func init() {
	fmt.Printf("github.com/gobuffalo/buffalo/middleware/csrf is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-csrf instead.")
}
