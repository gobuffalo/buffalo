package fizz

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Exec(t *testing.T) {
	r := require.New(t)

	b := NewBubbler(nil)
	f := fizzer{b}
	bb := &bytes.Buffer{}
	err := f.Exec(bb).(func(string) error)("echo hello")
	r.NoError(err)
	r.Equal("hello\n", bb.String())
}
