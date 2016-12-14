package resolvers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SimpleResolver(t *testing.T) {
	r := require.New(t)

	f, err := ioutil.TempFile(os.TempDir(), "example")
	r.NoError(err)
	defer os.Remove(f.Name())
	_, err = f.WriteString("hello")
	r.NoError(err)

	rr := &SimpleResolver{}
	b, err := rr.Read(f.Name())
	r.NoError(err)
	r.Equal("hello", string(b))

	_, err = rr.Read("unknown")
	r.Error(err)
}
