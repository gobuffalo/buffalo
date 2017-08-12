package plush

import (
	"testing"

	"golang.org/x/sync/errgroup"

	"github.com/stretchr/testify/require"
)

func Test_Context_Set(t *testing.T) {
	r := require.New(t)
	c := NewContext()
	r.Nil(c.Value("foo"))
	c.Set("foo", "bar")
	r.NotNil(c.Value("foo"))
}

func Test_Context_Set_Concurrency(t *testing.T) {
	r := require.New(t)
	c := NewContext()

	wg := errgroup.Group{}
	f := func() error {
		c.Set("a", "b")
		return nil
	}
	wg.Go(f)
	wg.Go(f)
	wg.Go(f)
	err := wg.Wait()
	r.NoError(err)
}

func Test_Context_Get(t *testing.T) {
	r := require.New(t)
	c := NewContext()
	r.Nil(c.Value("foo"))
	c.Set("foo", "bar")
	r.Equal("bar", c.Value("foo"))
}

func Test_NewSubContext_Set(t *testing.T) {
	r := require.New(t)

	c := NewContext()
	r.Nil(c.Value("foo"))

	sc := c.New()
	r.Nil(sc.Value("foo"))
	sc.Set("foo", "bar")
	r.Equal("bar", sc.Value("foo"))

	r.Nil(c.Value("foo"))
}

func Test_NewSubContext_Get(t *testing.T) {
	r := require.New(t)

	c := NewContext()
	c.Set("foo", "bar")

	sc := c.New()
	r.Equal("bar", sc.Value("foo"))
}
