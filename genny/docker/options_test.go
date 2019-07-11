package docker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		Style: "foo",
	}
	err := opts.Validate()
	r.Error(err)

	opts.Style = "multi"
	err = opts.Validate()
	r.NoError(err)

}
