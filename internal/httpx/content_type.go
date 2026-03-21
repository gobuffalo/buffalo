// Package httpx provides HTTP utility functions.
package httpx

import (
	"mime"
	"net/http"
	"strings"
)

// ContentType extracts the content type from the request.
// It checks Content-Type header first, then falls back to Accept header.
// Returns the media type without parameters (e.g., "application/json").
// Empty string is returned if only wildcard "*/*" is present.
func ContentType(req *http.Request) string {
	if ct := req.Header.Get("Content-Type"); ct != "" {
		mediatype, _, err := mime.ParseMediaType(ct)
		if err == nil && mediatype != "" {
			return mediatype
		}
	}

	accept := req.Header.Get("Accept")
	if accept == "" {
		return ""
	}

	for part := range strings.SplitSeq(accept, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		mediatype, _, err := mime.ParseMediaType(part)
		if err != nil {
			continue
		}

		if mediatype == "*/*" || mediatype == "*" {
			continue
		}

		return mediatype
	}

	return ""
}
