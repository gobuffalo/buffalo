package buffalo

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Response_MultipleWrite(t *testing.T) {
	r := require.New(t)
	resWr := httptest.NewRecorder()
	res := Response{
		ResponseWriter: resWr,
	}

	res.WriteHeader(http.StatusOK)
	res.WriteHeader(http.StatusInternalServerError)

	r.Equal(res.Status, http.StatusOK)
}
