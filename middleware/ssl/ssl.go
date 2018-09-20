package ssl

import (
	forcessl "github.com/gobuffalo/mw-forcessl"
	"github.com/markbates/oncer"
)

// ForceSSL uses the github.com/unrolled/secure package to
// automatically force a redirect to https from http.
// See https://github.com/unrolled/secure/ for more details
// on configuring.
//
// Deprecated: use github.com/gobuffalo/mw-forcessl#Middleware instead.
var ForceSSL = forcessl.Middleware

func init() {
	oncer.Deprecate(0, "github.com/gobuffalo/buffalo/middleware/ssl", "Use github.com/gobuffalo/mw-forcessl instead.")
}
