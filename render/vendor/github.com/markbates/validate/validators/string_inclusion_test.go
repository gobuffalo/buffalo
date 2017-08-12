package validators_test

import (
	"testing"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_StringInclusion(t *testing.T) {
	r := require.New(t)

	l := []string{"Mark", "Bates"}

	v := StringInclusion{"Name", "Mark", l}
	errors := validate.NewErrors()
	v.IsValid(errors)
	r.Equal(errors.Count(), 0)

	v = StringInclusion{"Name", "Foo", l}
	v.IsValid(errors)
	r.Equal(errors.Count(), 1)
	r.Equal(errors.Get("name"), []string{"Name is not in the list [Mark, Bates]."})
}
