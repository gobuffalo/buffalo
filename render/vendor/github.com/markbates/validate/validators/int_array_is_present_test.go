package validators_test

import (
	"testing"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_IntArrayIsPresent(t *testing.T) {
	r := require.New(t)

	v := IntArrayIsPresent{"Name", []int{1}}
	errors := validate.NewErrors()
	v.IsValid(errors)
	r.Equal(errors.Count(), 0)

	v = IntArrayIsPresent{"Name", []int{}}
	v.IsValid(errors)
	r.Equal(errors.Count(), 1)
	r.Equal(errors.Get("name"), []string{"Name can not be empty."})
}
