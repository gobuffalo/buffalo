package resource

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New_WithNestedName(t *testing.T) {
	r := require.New(t)

	g, err := New("", "admin/user")
	r.NoError(err)
	name := g.Name
	r.Equal("admin_user_id", name.ParamID())
}
