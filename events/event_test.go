package events

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Event_Validate(t *testing.T) {
	r := require.New(t)

	e := Event{}
	r.Error(e.Validate())

	e.Kind = "foo"
	r.NoError(e.Validate())
}
