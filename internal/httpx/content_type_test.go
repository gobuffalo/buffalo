package httpx

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ContentType(t *testing.T) {
	r := require.New(t)

	table := []struct {
		Header   string
		Value    string
		Expected string
	}{
		{"content-type", "a", "a"},
		{"Content-Type", "c,d", "c"},
		{"Content-Type", "e;f", "e"},
		{"Content-Type", "", ""},
		{"Accept", "", ""},
		{"Accept", "*/*", ""},
		{"Accept", "*/*;q=0.5, text/javascript, application/javascript, application/ecmascript, application/x-ecmascript", "text/javascript"},
		{"accept", "text/javascript,application/javascript,application/ecmascript,application/x-ecmascript", "text/javascript"},
	}

	for _, tt := range table {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(tt.Header, tt.Value)
		r.Equal(tt.Expected, ContentType(req))
	}
}
