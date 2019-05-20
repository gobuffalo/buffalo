package plugdeps

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Plugins_Encode(t *testing.T) {
	r := require.New(t)

	bb := &bytes.Buffer{}

	plugs := New()
	plugs.Add(pop, heroku, local)

	r.NoError(plugs.Encode(bb))

	fmt.Println(bb.String())
	act := strings.TrimSpace(bb.String())
	exp := strings.TrimSpace(eToml)
	r.Equal(exp, act)
}

func Test_Plugins_Decode(t *testing.T) {
	r := require.New(t)

	plugs := New()
	r.NoError(plugs.Decode(strings.NewReader(eToml)))

	names := []string{"buffalo-hello.rb", "buffalo-heroku", "buffalo-plugins", "buffalo-pop"}
	list := plugs.List()

	r.Len(list, len(names))
	for i, p := range list {
		r.Equal(names[i], p.Binary)
	}
}

func Test_Plugins_Remove(t *testing.T) {
	r := require.New(t)

	plugs := New()
	plugs.Add(pop, heroku)
	r.Len(plugs.List(), 3)
	plugs.Remove(pop)
	r.Len(plugs.List(), 2)
}
