package core

import (
	"go/build"
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_validateInGoPath(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	c := build.Default
	err := validateInGoPath(c.SrcDirs())(run)
	r.NoError(err)

	err = validateInGoPath([]string{"idontexist"})(run)
	r.Error(err)
}
