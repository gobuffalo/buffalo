package httpx

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ContentType(t *testing.T) {
	r := require.New(t)

	table := []struct {
		name     string
		header   string
		value    string
		expected string
	}{
		{"simple content-type", "content-type", "application/json", "application/json"},
		{"content-type with charset", "Content-Type", "application/json; charset=utf-8", "application/json"},
		{"accept single value", "Accept", "text/html", "text/html"},
		{"accept with quality", "Accept", "text/html;q=0.9", "text/html"},
		{"accept multiple values", "Accept", "application/json, text/html", "application/json"},
		{"accept skips wildcard", "Accept", "*/*, application/json", "application/json"},
		{"accept with quality values", "Accept", "text/plain;q=0.5, application/json;q=0.9", "text/plain"},
		{"empty content-type", "Content-Type", "", ""},
		{"empty accept", "Accept", "", ""},
		{"wildcard only", "Accept", "*/*", ""},
		{"complex accept header", "Accept", "*/*;q=0.5, text/javascript, application/javascript", "text/javascript"},
		{"lowercase header name", "accept", "application/json", "application/json"},
		{"uppercase media type", "Accept", "Application/JSON", "application/json"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set(tt.header, tt.value)
			r.Equal(tt.expected, ContentType(req))
		})
	}
}

func Test_ContentType_Priority(t *testing.T) {
	r := require.New(t)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/html")
	r.Equal("application/json", ContentType(req))

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "text/html")
	r.Equal("text/html", ContentType(req))
}
