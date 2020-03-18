// +build integration_test

package integration

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gobuffalo/pop/v5"
	"github.com/stretchr/testify/require"
)

func Test_NewCmd_NoName(t *testing.T) {
	err := call([]string{"new"}, nil)
	r := require.New(t)
	r.Error(err)
	r.EqualError(err, "you must enter a name for your new application")
}

func Test_NewCmd_InvalidDBType(t *testing.T) {
	args := []string{
		"new",
		"coke",
		"--db-type",
		"a",
	}
	err := call(args, nil)
	r := require.New(t)
	r.Error(err)
	r.EqualError(err, fmt.Sprintf(`unknown dialect "a" expecting one of %s`, strings.Join(pop.AvailableDialects, ", ")))
}

func Test_NewCmd_ForbiddenAppName(t *testing.T) {
	args := []string{
		"new",
		"buffalo",
	}
	err := call(args, nil)
	r := require.New(t)
	r.Error(err)
	r.EqualError(err, "name buffalo is not allowed, try a different application name")
}
