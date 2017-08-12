package packr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Box_String(t *testing.T) {
	r := require.New(t)
	s := testBox.String("hello.txt")
	r.Equal("hello world!\n", s)
}

func Test_Box_MustString(t *testing.T) {
	r := require.New(t)
	_, err := testBox.MustString("idontexist.txt")
	r.Error(err)
}

func Test_Box_Bytes(t *testing.T) {
	r := require.New(t)
	s := testBox.Bytes("hello.txt")
	r.Equal([]byte("hello world!\n"), s)
}

func Test_Box_MustBytes(t *testing.T) {
	r := require.New(t)
	_, err := testBox.MustBytes("idontexist.txt")
	r.Error(err)
}

func Test_Box_Walk_Physical(t *testing.T) {
	r := require.New(t)
	count := 0
	err := testBox.Walk(func(path string, f File) error {
		count++
		return nil
	})
	r.NoError(err)
	r.Equal(2, count)
}

func Test_Box_Walk_Virtual(t *testing.T) {
	r := require.New(t)
	count := 0
	err := virtualBox.Walk(func(path string, f File) error {
		count++
		return nil
	})
	r.NoError(err)
	r.Equal(3, count)
}
