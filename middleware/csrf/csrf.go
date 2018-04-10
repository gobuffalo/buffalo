package csrf

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/markbates/going/defaults"
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
	ErrNoReferer = errors.New("referer not supplied")
	// ErrBadReferer is returned when the scheme & host in the URL do not match
	// the supplied Referer header.
	ErrBadReferer = errors.New("referer invalid")
	// ErrNoToken is returned if no CSRF token is supplied in the request.
	ErrNoToken = errors.New("CSRF token not found in request")
	// ErrBadToken is returned if the CSRF token in the request does not match
	// the token in the session, or is otherwise malformed.
	ErrBadToken = errors.New("CSRF token invalid")
)

// New enable CSRF protection on routes using this middleware.
// This middleware is adapted from gorilla/csrf
var New = func(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		// don't run in test mode
		if envy.Get("GO_ENV", "development") == "test" {
			c.Set(tokenKey, "test")
			return next(c)
		}

		req := c.Request()

		ct := defaults.String(req.Header.Get("Content-Type"), req.Header.Get("Accept"))
		// ignore non-html requests
		if ct != "" && !contains(htmlTypes, ct) {
			return next(c)
		}

		var realToken []byte
		rawRealToken := c.Session().Get(tokenKey)

		if rawRealToken == nil || len(rawRealToken.([]byte)) != tokenLength {
			// If the token is missing, or the length if the token is wrong,
			// generate a new token.
			realToken, err := generateRandomBytes(tokenLength)
			if err != nil {
				return err
			}
			// Save the new real token in session
			c.Session().Set(tokenKey, realToken)
		} else {
			realToken = rawRealToken.([]byte)
		}

		// Set masked token in context data, to be available in template
		c.Set(fieldName, mask(realToken, req))

		// HTTP methods not defined as idempotent ("safe") under RFC7231 require
		// inspection.
		if !contains(safeMethods, req.Method) {
			// Enforce an origin check for HTTPS connections. As per the Django CSRF
			// implementation (https://goo.gl/vKA7GE) the Referer header is almost
			// always present for same-domain HTTP requests.
			if req.URL.Scheme == "https" {
				// Fetch the Referer value. Call the error handler if it's empty or
				// otherwise fails to parse.
				referer, err := url.Parse(req.Referer())
				if err != nil || referer.String() == "" {
					return ErrNoReferer
				}

				if !sameOrigin(req.URL, referer) {
					return ErrBadReferer
				}
			}

			// Retrieve the combined token (pad + masked) token and unmask it.
			requestToken := unmask(requestCSRFToken(req))

			// Missing token
			if requestToken == nil {
				return ErrNoToken
			}

			// Compare tokens
			if !compareTokens(requestToken, realToken) {
				return ErrBadToken
			}
		}

		return next(c)
	}
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random number generator
// fails to function correctly.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// err == nil only if len(b) == n
	if err != nil {
		return nil, err
	}

	return b, nil
}

// sameOrigin returns true if URLs a and b share the same origin. The same
// origin is defined as host (which includes the port) and scheme.
func sameOrigin(a, b *url.URL) bool {
	return (a.Scheme == b.Scheme && a.Host == b.Host)
}

// contains is a helper function to check if a string exists in a slice - e.g.
// whether a HTTP method exists in a list of safe methods.
func contains(vals []string, s string) bool {
	s = strings.ToLower(s)
	for _, v := range vals {
		if strings.Contains(s, strings.ToLower(v)) {
			return true
		}
	}

	return false
}

// compare securely (constant-time) compares the unmasked token from the request
// against the real token from the session.
func compareTokens(a, b []byte) bool {
	// This is required as subtle.ConstantTimeCompare does not check for equal
	// lengths in Go versions prior to 1.3.
	if len(a) != len(b) {
		return false
	}

	return subtle.ConstantTimeCompare(a, b) == 1
}

// xorToken XORs tokens ([]byte) to provide unique-per-request CSRF tokens. It
// will return a masked token if the base token is XOR'ed with a one-time-pad.
// An unmasked token will be returned if a masked token is XOR'ed with the
// one-time-pad used to mask it.
func xorToken(a, b []byte) []byte {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}

	res := make([]byte, n)

	for i := 0; i < n; i++ {
		res[i] = a[i] ^ b[i]
	}

	return res
}

// mask returns a unique-per-request token to mitigate the BREACH attack
// as per http://breachattack.com/#mitigations
//
// The token is generated by XOR'ing a one-time-pad and the base (session) CSRF
// token and returning them together as a 64-byte slice. This effectively
// randomises the token on a per-request basis without breaking multiple browser
// tabs/windows.
func mask(realToken []byte, r *http.Request) string {
	otp, err := generateRandomBytes(tokenLength)
	if err != nil {
		return ""
	}

	// XOR the OTP with the real token to generate a masked token. Append the
	// OTP to the front of the masked token to allow unmasking in the subsequent
	// request.
	return base64.StdEncoding.EncodeToString(append(otp, xorToken(otp, realToken)...))
}

// unmask splits the issued token (one-time-pad + masked token) and returns the
// unmasked request token for comparison.
func unmask(issued []byte) []byte {
	// Issued tokens are always masked and combined with the pad.
	if len(issued) != tokenLength*2 {
		return nil
	}

	// We now know the length of the byte slice.
	otp := issued[tokenLength:]
	masked := issued[:tokenLength]

	// Unmask the token by XOR'ing it against the OTP used to mask it.
	return xorToken(otp, masked)
}

// requestCSRFToken gets the CSRF token from either:
// - a HTTP header
// - a form value
// - a multipart form value
func requestCSRFToken(r *http.Request) []byte {
	// 1. Check the HTTP header first.
	issued := r.Header.Get(headerName)

	// 2. Fall back to the POST (form) value.
	if issued == "" {
		issued = r.PostFormValue(fieldName)
	}

	// 3. Finally, fall back to the multipart form (if set).
	if issued == "" && r.MultipartForm != nil {
		vals := r.MultipartForm.Value[fieldName]

		if len(vals) > 0 {
			issued = vals[0]
		}
	}

	// Decode the "issued" (pad + masked) token sent in the request. Return a
	// nil byte slice on a decoding error (this will fail upstream).
	decoded, err := base64.StdEncoding.DecodeString(issued)
	if err != nil {
		return nil
	}

	return decoded
}
