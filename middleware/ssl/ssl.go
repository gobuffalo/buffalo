package ssl

import (
	"fmt"

	forcessl "github.com/gobuffalo/mw-forcessl"
)

// ForceSSL uses the github.com/unrolled/secure package to
// automatically force a redirect to https from http.
// See https://github.com/unrolled/secure/ for more details
// on configuring.
//
// Deprecated: use github.com/gobuffalo/mw-forcessl#Middleware instead.
var ForceSSL = forcessl.Middleware

func init() {
	fmt.Printf("github.com/gobuffalo/buffalo/middleware/ssl is deprecated and will be removed in the next version. Please use github.com/gobuffalo/mw-forcessl instead.")
}
