package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	a := assert.New(t)

	a.Equal(String("", "foo"), "foo")
	a.Equal(String("bar", "foo"), "bar")
	var s string
	a.Equal(String(s, "foo"), "foo")
}

func Test_Int(t *testing.T) {
	a := assert.New(t)

	a.Equal(Int(0, 1), 1)
	a.Equal(Int(2, 1), 2)
	var s int
	a.Equal(Int(s, 1), 1)
}

func Test_Int64(t *testing.T) {
	a := assert.New(t)

	a.Equal(Int64(0, 1), int64(1))
	a.Equal(Int64(2, 1), int64(2))
	var s int64
	a.Equal(Int64(s, 1), int64(1))
}

func Test_Float32(t *testing.T) {
	a := assert.New(t)

	a.Equal(Float32(0, 1), float32(1))
	a.Equal(Float32(2, 1), float32(2))
	var s float32
	a.Equal(Float32(s, 1), float32(1))
}

func Test_Float64(t *testing.T) {
	a := assert.New(t)

	a.Equal(Float64(0, 1), float64(1))
	a.Equal(Float64(2, 1), float64(2))
	var s float64
	a.Equal(Float64(s, 1), float64(1))
}
