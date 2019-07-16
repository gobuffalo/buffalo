package render

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Func(t *testing.T) {
	r := require.New(t)

	table := []rendFriend{
		Func,
		New(Options{}).Func,
	}

	for _, tt := range table {
		bb := &bytes.Buffer{}

		re := tt("foo/bar", func(w io.Writer, data Data) error {
			_, err := w.Write([]byte(data["name"].(string)))
			return err
		})

		r.Equal("foo/bar", re.ContentType())
		err := re.Render(bb, Data{"name": "Mark"})
		r.NoError(err)
		r.Equal("Mark", bb.String())
	}

}
