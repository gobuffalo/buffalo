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

func Test_New_WithPropertyNames(t *testing.T) {
	r := require.New(t)

	cases := []struct {
		Case     string
		Args     []string
		HasError bool
	}{
		{Case: "", Args: []string{"", "body,name"}, HasError: true},
		{Case: "", Args: []string{"", "body,name:nulls.String"}, HasError: true},
		{Case: "", Args: []string{"", "body", "name:nulls.String"}, HasError: false},
		{Case: "", Args: []string{"", "body.name"}, HasError: true},
		{Case: "", Args: []string{"", "body", "name"}, HasError: false},
		{Case: "", Args: []string{"", "body-name"}, HasError: false},
		{Case: "", Args: []string{"", "body_name"}, HasError: false},
		{Case: "", Args: []string{"", "9A"}, HasError: true},
		{Case: "", Args: []string{"", "AAA"}, HasError: false},
	}

	for _, tcase := range cases {
		_, err := New("", tcase.Args...)

		if tcase.HasError {
			r.Error(err, tcase.Case)
			continue
		}

		r.NoError(err, tcase.Case)
	}
}
