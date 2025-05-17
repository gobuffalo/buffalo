package buffalo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newRequestData(t *testing.T) {

	r := require.New(t)
	ts := newRequestData()
	r.NotNil(ts)
	r.NotNil(ts.moot)
	r.NotNil(ts.d)
}
