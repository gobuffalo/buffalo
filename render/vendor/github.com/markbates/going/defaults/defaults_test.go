package defaults_test

import (
	"testing"

	"github.com/markbates/going/defaults"
	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	a := assert.New(t)

	a.Equal(defaults.String("", "foo"), "foo")
	a.Equal(defaults.String("bar", "foo"), "bar")
	var s string
	a.Equal(defaults.String(s, "foo"), "foo")
}

func Test_Int(t *testing.T) {
	a := assert.New(t)

	a.Equal(defaults.Int(0, 1), 1)
	a.Equal(defaults.Int(2, 1), 2)
	var s int
	a.Equal(defaults.Int(s, 1), 1)
}

func Test_Int64(t *testing.T) {
	a := assert.New(t)

	a.Equal(defaults.Int64(0, 1), int64(1))
	a.Equal(defaults.Int64(2, 1), int64(2))
	var s int64
	a.Equal(defaults.Int64(s, 1), int64(1))
}

func Test_Float32(t *testing.T) {
	a := assert.New(t)

	a.Equal(defaults.Float32(0, 1), float32(1))
	a.Equal(defaults.Float32(2, 1), float32(2))
	var s float32
	a.Equal(defaults.Float32(s, 1), float32(1))
}

func Test_Float64(t *testing.T) {
	a := assert.New(t)

	a.Equal(defaults.Float64(0, 1), float64(1))
	a.Equal(defaults.Float64(2, 1), float64(2))
	var s float64
	a.Equal(defaults.Float64(s, 1), float64(1))
}
