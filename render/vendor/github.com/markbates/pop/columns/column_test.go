package columns_test

import (
	"testing"

	"github.com/markbates/pop/columns"
	"github.com/stretchr/testify/require"
)

func Test_Column_UpdateString(t *testing.T) {
	r := require.New(t)
	c := columns.Column{Name: "foo"}
	r.Equal(c.UpdateString(), "foo = :foo")
}
