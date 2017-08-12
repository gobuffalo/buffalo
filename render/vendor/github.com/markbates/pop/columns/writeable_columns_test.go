package columns_test

import (
	"testing"

	"github.com/markbates/pop/columns"
	"github.com/stretchr/testify/require"
)

func Test_Columns_WriteableString_Symbolized(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ColumnsForStruct(f, "foo")
		u := c.Writeable().SymbolizedString()
		r.Equal(u, ":LastName, :write")
	}
}

func Test_Columns_UpdateString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ColumnsForStruct(f, "foo")
		u := c.Writeable().UpdateString()
		r.Equal(u, "LastName = :LastName, write = :write")
	}
}

func Test_Columns_WriteableString(t *testing.T) {
	r := require.New(t)
	for _, f := range []interface{}{foo{}, &foo{}} {
		c := columns.ColumnsForStruct(f, "foo")
		u := c.Writeable().String()
		r.Equal(u, "LastName, write")
	}
}
