package buffalo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FlashSet(t *testing.T) {
	r := require.New(t)
	f := newFlash(&Session{})

	r.Equal(f.data, map[string][]string{})
	f.Set("error", "something")
	r.Equal(f.data, map[string][]string{
		"error": []string{"something"},
	})
}

func Test_FlashGet(t *testing.T) {
	r := require.New(t)
	f := newFlash(&Session{})

	r.Equal(f.data, map[string][]string{})
	f.Set("error", "something")
	r.Equal(f.Get("error"), []string{"something"})
}

func Test_FlashDelete(t *testing.T) {
	r := require.New(t)
	f := newFlash(&Session{})

	r.Equal(f.data, map[string][]string{})
	f.Set("error", "something")
	r.Equal(f.Get("error"), []string{"something"})

	f.Delete("error")
	r.Equal(f.Get("error"), []string(nil))
}

func Test_FlashClear(t *testing.T) {
	r := require.New(t)
	f := newFlash(&Session{})

	r.Equal(f.data, map[string][]string{})
	f.Set("error", "something")
	f.Set("warning", "warning")
	r.Equal(f.Get("error"), []string{"something"})
	r.Equal(f.Get("warning"), []string{"warning"})

	f.Clear()
	r.Equal(f.data, map[string][]string{})

	r.Equal(f.Get("error"), []string(nil))
	r.Equal(f.Get("warning"), []string(nil))
}

func Test_FlashAdd(t *testing.T) {
	r := require.New(t)
	f := newFlash(&Session{})

	r.Equal(f.data, map[string][]string{})

	f.Add("error", "something")
	r.Equal(f.data, map[string][]string{
		"error": []string{"something"},
	})

	f.Add("error", "other")
	r.Equal(f.data, map[string][]string{
		"error": []string{"something", "other"},
	})
}
