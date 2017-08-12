package validators_test

import (
	"testing"

	"github.com/markbates/validate"
	. "github.com/markbates/validate/validators"
	"github.com/stretchr/testify/require"
)

func Test_IntIsPresent(t *testing.T) {
	r := require.New(t)

	v := IntIsPresent{"Name", 1}
	errors := validate.NewErrors()
	v.IsValid(errors)
	r.Equal(errors.Count(), 0)

	v = IntIsPresent{"Name", 0}
	v.IsValid(errors)
	r.Equal(errors.Count(), 1)
	r.Equal(errors.Get("name"), []string{"Name can not be blank."})
}
