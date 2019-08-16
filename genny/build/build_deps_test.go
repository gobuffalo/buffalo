package build

import (
	"strings"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_buildDeps(t *testing.T) {
	r := require.New(t)

	opts := &Options{
		Tags: meta.BuildTags{"foo"},
	}

	run := gentest.NewRunner()
	run.WithNew(buildDeps(opts))

	r.NoError(run.Run())

	res := run.Results()

	if envy.Mods() {
		r.Len(res.Commands, 0)
		return
	}
	r.Len(res.Commands, 1)

	c := res.Commands[0]
	r.Equal("go get -tags development foo ./...", strings.Join(c.Args, " "))
}
