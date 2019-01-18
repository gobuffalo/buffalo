package actions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	err := opts.Validate()
	r.Error(err)

	opts.Name = "user"

	err = opts.Validate()
	r.Error(err)

	opts.Actions = []string{"index", "show"}
	err = opts.Validate()
	r.NoError(err)
}
