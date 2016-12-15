package resolvers

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GoPathResolver(t *testing.T) {
	r := require.New(t)
	oldpath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldpath)
	os.Setenv("GOPATH", os.TempDir())

	gp := filepath.Join(os.Getenv("GOPATH"), "src", "foo", "bar")
	os.MkdirAll(gp, 0755)
	f, err := ioutil.TempFile(gp, "example")
	r.NoError(err)
	defer os.Remove(gp)
	_, err = f.WriteString("hello")
	r.NoError(err)

	rr := &GoPathResolver{}
	b, err := rr.Read(filepath.Base(f.Name()))
	r.NoError(err)
	r.Equal("hello", string(b))

	_, err = rr.Read("unknown")
	r.Error(err)
}
