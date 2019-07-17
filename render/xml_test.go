package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_XML(t *testing.T) {
	r := require.New(t)

	type user struct {
		Name string
	}

	e := NewEngine()

	re := e.XML(user{Name: "Mark"})
	r.Equal("application/xml; charset=utf-8", re.ContentType())

	bb := &bytes.Buffer{}

	r.NoError(re.Render(bb, nil))
	r.Equal("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<user>\n  <Name>Mark</Name>\n</user>", strings.TrimSpace(bb.String()))
}
