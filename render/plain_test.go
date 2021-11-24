package render

import (
	"bytes"
	"testing"

	"github.com/psanford/memfs"
	"github.com/stretchr/testify/require"
)

func Test_Plain(t *testing.T) {
	r := require.New(t)

	rootFS := memfs.New()
	r.NoError(rootFS.WriteFile("test.txt", []byte("<%= name %>"), 0644))

	e := NewEngine()
	e.TemplatesFS = rootFS

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
