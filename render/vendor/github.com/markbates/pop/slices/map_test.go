package slices

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Map_UnmarshalText(t *testing.T) {
	r := require.New(t)

	m := Map{}
	err := m.UnmarshalText([]byte(`{"a":"b"}`))
	r.NoError(err)
	r.Equal("b", m["a"])
}

func Test_Map_MarshalJSON(t *testing.T) {
	r := require.New(t)

	m := Map{"a": "b"}
	b, err := json.Marshal(m)
	r.NoError(err)
	r.Equal([]byte(`{"a":"b"}`), b)
}
