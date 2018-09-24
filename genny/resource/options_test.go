package resource

import (
	"testing"

	"github.com/gobuffalo/flect/name"
	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{}

	r.Error(opts.Validate())

	opts.Name = name.New("widget")
	r.NoError(opts.Validate())

	opts = &Options{
		Args: []string{"widget"},
	}
	r.NoError(opts.Validate())
	r.Equal("widget", opts.Name.String())

	opts = &Options{
		Args: []string{"widget", "name", "bio:nulls.Text"},
	}
	r.NoError(opts.Validate())
	r.Equal("widget", opts.Name.String())
	r.Len(opts.Attrs, 2)

	a := opts.Attrs[0]
	r.Equal("name", a.Name.String())
	r.Equal("string", a.Type)

	a = opts.Attrs[1]
	r.Equal("bio", a.Name.String())
	r.Equal("text", a.Type)
}
