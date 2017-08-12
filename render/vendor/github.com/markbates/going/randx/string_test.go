package randx_test

import (
	"math/rand"
	"testing"

	"github.com/markbates/going/randx"
	"github.com/stretchr/testify/require"
)

func init() {
	rand.Seed(1)
}

func Test_String(t *testing.T) {
	r := require.New(t)
	r.Len(randx.String(5), 5)
	r.Len(randx.String(50), 50)
}
