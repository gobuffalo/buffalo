package render

import (
	"net/http"
	"strings"
)

// FromContentType renders the value using the content type
// detected in "Content-Type" header. If the header doesn't exist
// it tries to fetch it from the "format" URL query argument.
// However, if it fails to detect the content type, json is provided
// as a fallback.
func FromContentType(v interface{}, req *http.Request) Renderer {
	if contentType, found := req.Header["Content-Type"]; found {
		// Try to read content type from Content-Type HTTP header
		if strings.EqualFold(contentType[0], "application/json") {
			return jsonRenderer{value: v}
		} else if strings.EqualFold(contentType[0], "application/xml") {
			return xmlRenderer{value: v}
		}
	} else {
		// Try to get content type as a query argument
		format := req.URL.Query().Get("format")
		if len(format) > 0 {
			if strings.EqualFold(format, "json") {
				return jsonRenderer{value: v}
			} else if strings.EqualFold(format, "xml") {
				return xmlRenderer{value: v}
			}
		}
	}

	// jsonRenderer as fallback
	return jsonRenderer{value: v}
}

// FromContentType renders the value using the content type
// detected in "Content-Type" header. If the header doesn't exist
// it tries to fetch it from the "format" URL query argument.
// However, if it fails to detect the content type, json is provided
// as a fallback.
func (e *Engine) FromContentType(v interface{}, req *http.Request) Renderer {
	return FromContentType(v, req)
}
