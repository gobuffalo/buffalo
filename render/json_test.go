package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_JSON(t *testing.T) {
	r := require.New(t)

	e := NewEngine()

	re := e.JSON(Data{"hello": "world"})
	r.Equal("application/json; charset=utf-8", re.ContentType())

	bb := &bytes.Buffer{}

	r.NoError(re.Render(bb, nil))
	r.Equal(`{"hello":"world"}`, strings.TrimSpace(bb.String()))
}
