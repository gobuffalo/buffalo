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
	r.Equal("admin_user_id", name.ParamID().String())
}

func Test_New_WithNestedName_Plural(t *testing.T) {
	r := require.New(t)

	g, err := New("", "admin/users")
	r.NoError(err)
	name := g.Name
	r.Equal("admin_user_id", name.ParamID().String())
}

func Test_New_FilesPath_WithNestedName(t *testing.T) {
	r := require.New(t)

	g, err := New("", "Admin/superFast/Plane")
	r.NoError(err)
	r.Equal("admin/super_fast/planes", g.FilesPath)
	r.Equal("admin_super_fast_planes", g.ActionsPath)
}
