package newapp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_nameHasIllegalCharacter(t *testing.T) {
	m := map[string]bool{
		"coke":                      false,
		"my-coke":                   false,
		"my_coke":                   false,
		"COKE":                      false,
		"MY-COKE":                   false,
		"MY_COKE":                   false,
		"123COKE":                   false,
		"1(3c&ke":                   true,
		"github.com/markbates/coke": true,
	}
	for k, v := range m {
		g, _ := New(k)
		t.Run(k, func(st *testing.T) {
			r := require.New(st)
			if v {
				r.Error(g.Validate())
			} else {
				r.NoError(g.Validate())
			}
		})
	}
}
