package build

import (
	"sync"
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/meta"
	"github.com/stretchr/testify/require"
)

func Test_WithDeps(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()

	opts := &Options{
		WithAssets:  false,
		Environment: "bar",
		App:         meta.New("."),
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
