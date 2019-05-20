package plugcmds

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gobuffalo/events"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func Test_Available_Add(t *testing.T) {
	r := require.New(t)

	a := NewAvailable()
	err := a.Add("generate", &cobra.Command{
		Use:     "foo",
		Short:   "generates foo",
		Aliases: []string{"f"},
	})
	r.NoError(err)

	r.Len(a.Commands(), 2)
}

func Test_Available_Encode(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	a := NewAvailable()
	err := a.Add("generate", &cobra.Command{
		Use:     "foo",
		Short:   "generates foo",
		Aliases: []string{"f"},
	})
	r.NoError(err)

	r.NoError(a.Encode(bb))
	const exp = `[{"name":"foo","use_command":"foo","buffalo_command":"generate","description":"generates foo","aliases":["f"]}]`
	r.Equal(exp, strings.TrimSpace(bb.String()))
}

func Test_Available_Listen(t *testing.T) {
	r := require.New(t)

	a := NewAvailable()
	err := a.Listen(func(e events.Event) error {
		return nil
	})
	r.NoError(err)

	r.Len(a.Commands(), 2)
}
