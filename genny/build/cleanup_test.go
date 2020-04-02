package build

import (
	"sync"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/v2/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_WithDeps(t *testing.T) {
	r := require.New(t)
	envy.Set(envy.GO111MODULE, "on")

	run := gentest.NewRunner()

	opts := &Options{
		WithAssets:    false,
		WithBuildDeps: true,
		Environment:   "bar",
		App:           meta.New("."),
	}

	emptyMap := sync.Map{}
	opts.rollback = &emptyMap

	f := Cleanup(opts)
	f(run)

	results := run.Results()

	cmds := []string{"go mod tidy"}
	for i, c := range results.Commands {
		eq(r, cmds[i], c)
	}
}

func Test_WithoutDeps(t *testing.T) {
	r := require.New(t)
	envy.Set(envy.GO111MODULE, "on")

	run := gentest.NewRunner()

	opts := &Options{
		WithAssets:    false,
		WithBuildDeps: false,
		Environment:   "bar",
		App:           meta.New("."),
	}

	emptyMap := sync.Map{}
	opts.rollback = &emptyMap

	f := Cleanup(opts)
	f(run)

	results := run.Results()

	r.Len(results.Commands, 0)
}
