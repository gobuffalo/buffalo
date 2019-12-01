package buffalo

import "net/http"

// PreWare takes an http.Handler and returns and http.Handler
// and acts as a pseudo-middleware between the http.Server and
// a Buffalo application.
type PreWare func(http.Handler) http.Handler
