package info

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	err := opts.Validate()
	r.Error(err)

  // TODO: make it valid. :)

	err = opts.Validate()
  r.NoError(err)

  // TODO: add some more assertions here.
}
