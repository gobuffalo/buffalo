package ci

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	err := opts.Validate()
	r.Error(err)

	opts.Provider = "travis-ci"
	opts.DBType = "postgres"

	err = opts.Validate()
	r.NoError(err)
}
