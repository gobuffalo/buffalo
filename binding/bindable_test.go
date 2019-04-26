package binding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type orbison struct {
	bound bool
}

func (o *orbison) Bind(req *http.Request) error {
	o.bound = true
	return nil
}

func Test_Bindable(t *testing.T) {
	r := require.New(t)

	req := httptest.NewRequest("GET", "/", nil)
	o := &orbison{}
	r.False(o.bound)
	r.NoError(Exec(req, o))
	r.True(o.bound)
}
