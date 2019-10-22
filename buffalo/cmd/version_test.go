package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_VersionCmd(t *testing.T) {
	r := require.New(t)

	c := RootCmd
	c.SetArgs([]string{
		"version",
	})
	err := c.Execute()
	r.NoError(err)
}

func Test_VersionJSONCmd(t *testing.T) {
	r := require.New(t)

	c := RootCmd
	c.SetArgs([]string{
		"version",
		"--json",
	})
	err := c.Execute()
	r.NoError(err)
}
