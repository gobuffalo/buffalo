package validators_test

import (
	"testing"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_IntIsGreaterThan(t *testing.T) {
	r := require.New(t)

	v := IntIsGreaterThan{Name: "Number", Field: 2, Compared: 1}
	errors := validate.NewErrors()
	v.IsValid(errors)
	r.Equal(0, errors.Count())

	v = IntIsGreaterThan{Name: "number", Field: 1, Compared: 2}
	v.IsValid(errors)
	r.Equal(1, errors.Count())
	r.Equal(errors.Get("number"), []string{"1 is not greater than 2."})
}
