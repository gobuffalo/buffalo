package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Bootstrap4_Default(t *testing.T) {
	r := require.New(t)
	f, err := newCmd.Flags().GetInt("bootstrap")
	r.NoError(err)
	r.Equal(4, f)
}
