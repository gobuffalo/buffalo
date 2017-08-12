package columns_test

import (
	"testing"

	"github.com/markbates/pop/columns"
	"github.com/stretchr/testify/require"
)

func Test_Columns_ReadableString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ColumnsForStruct(f, "foo")
		u := c.Readable().String()
		r.Equal(u, "LastName, first_name, read")
	}
}

func Test_Columns_Readable_SelectString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ColumnsForStruct(f, "foo")
		u := c.Readable().SelectString()
		r.Equal(u, "first_name as f, foo.LastName, foo.read")
	}
}

func Test_Columns_ReadableString_Symbolized(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ColumnsForStruct(f, "foo")
		u := c.Readable().SymbolizedString()
		r.Equal(u, ":LastName, :first_name, :read")
	}
}
