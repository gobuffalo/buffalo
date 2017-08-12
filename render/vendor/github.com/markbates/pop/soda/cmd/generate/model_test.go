package generate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_newName(t *testing.T) {
	r := require.New(t)
	n := newName("carrot")
	r.Equal(n.File, "carrot")
	r.Equal(n.Proper, "Carrot")
	r.Equal(n.Table, "carrots")
	r.Equal(n.Plural, "Carrots")
	r.Equal(n.Char, "c")
}

func Test_newName_Plural(t *testing.T) {
	r := require.New(t)
	n := newName("carrots")
	r.Equal(n.File, "carrot")
	r.Equal("Carrots", n.Proper)
	r.Equal(n.Table, "carrots")
	r.Equal(n.Plural, "Carrots")
	r.Equal(n.Char, "c")
}

func Test_newName_multipleCamelCase(t *testing.T) {
	r := require.New(t)
	n := newName("carrotCake")
	r.Equal(n.File, "carrot_cake")
	r.Equal(n.Proper, "CarrotCake")
	r.Equal(n.Table, "carrot_cakes")
	r.Equal(n.Plural, "CarrotCakes")
	r.Equal(n.Char, "c")
}

func Test_newName_multipleCamelCase_Plural(t *testing.T) {
	r := require.New(t)
	n := newName("carrotCakes")
	r.Equal(n.File, "carrot_cake")
	r.Equal(n.Proper, "CarrotCakes")
	r.Equal(n.Table, "carrot_cakes")
	r.Equal(n.Plural, "CarrotCakes")
	r.Equal(n.Char, "c")
}

func Test_newName_multipleSnake(t *testing.T) {
	r := require.New(t)
	n := newName("carrot_cake")
	r.Equal(n.File, "carrot_cake")
	r.Equal(n.Proper, "CarrotCake")
	r.Equal(n.Table, "carrot_cakes")
	r.Equal(n.Plural, "CarrotCakes")
	r.Equal(n.Char, "c")
}

func Test_newName_multipleSnake_Plural(t *testing.T) {
	r := require.New(t)
	n := newName("carrot_cakes")
	r.Equal(n.File, "carrot_cake")
	r.Equal(n.Proper, "CarrotCakes")
	r.Equal(n.Table, "carrot_cakes")
	r.Equal(n.Plural, "CarrotCakes")
	r.Equal(n.Char, "c")
}
