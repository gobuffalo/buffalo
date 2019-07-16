package render

import (
	"bytes"
	"testing"

	"github.com/gobuffalo/packd"
	"github.com/stretchr/testify/require"
)

func Test_Plain(t *testing.T) {
	r := require.New(t)

	box := packd.NewMemoryBox()
	r.NoError(box.AddString("test.txt", "<%= name %>"))

	e := NewEngine()
	e.TemplatesBox = box

	re := e.Plain("test.txt")
	r.Equal("text/plain; charset=utf-8", re.ContentType())

	var examples = []string{"Mark", "Jém"}
	for _, example := range examples {
		example := example
		bb := &bytes.Buffer{}
		r.NoError(re.Render(bb, Data{"name": example}))
		r.Equal(example, bb.String())
	}
}
