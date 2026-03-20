package render

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func Test_Plain(t *testing.T) {
	r := require.New(t)

	rootFS := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte("<%= name %>"),
			Mode: 0644,
		},
	}

	e := NewEngine()
	e.TemplatesFS = rootFS

	re := e.Plain("test.txt")
	r.Equal("text/plain; charset=utf-8", re.ContentType())

	var examples = []string{"Mark", "Jém"}
	for _, example := range examples {
		bb := &bytes.Buffer{}
		r.NoError(re.Render(bb, Data{"name": example}))
		r.Equal(example, bb.String())
	}
}
