package render

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_stringMap(t *testing.T) {
	r := require.New(t)

	sm := &stringMap{}

	sm.Store("a", `A`)

	s, ok := sm.Load("a")
	r.True(ok)
	r.Equal(`A`, s)

	s, ok = sm.LoadOrStore("b", `B`)
	r.True(ok)
	r.Equal(`B`, s)

	s, ok = sm.LoadOrStore("b", `BB`)
	r.True(ok)
	r.Equal(`B`, s)

	var keys []string

	sm.Range(func(key string, value string) bool {
		keys = append(keys, key)
		return true
	})

	sort.Strings(keys)

	r.Equal(sm.Keys(), keys)

	sm.Delete("b")
	r.Equal([]string{"a", "b"}, keys)

	sm.Delete("b")
	_, ok = sm.Load("b")
	r.False(ok)

	func(m *stringMap) {
		m.Store("c", `C`)
	}(sm)
	s, ok = sm.Load("c")
	r.True(ok)
	r.Equal(`C`, s)
}
